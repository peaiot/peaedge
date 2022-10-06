package modbus_task

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/golimit"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/models"
)

const (
	ERRORVALUE          = 999999999999
	ConstProtoModbusRTU = "modbusrtu"
	ConstProtoModbusTCP = "modbustcp"
)

func StartModbusReadTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	app.ResetModbusDevConnStatus()
	var limitPool = golimit.NewGoLimitPool()
	for {
		var devices = make([]models.ModbusDevice, 0)
		err := app.DB().Find(&devices).Error
		if err != nil {
			emsg := fmt.Sprintf("读取设备列表失败... %s", err.Error())
			log.Errorf(emsg)
			time.Sleep(time.Second * 5)
			continue
		}

		wg := sync.WaitGroup{}
		wg.Add(len(devices))
		for _, dev := range devices {
			go func(devitem models.ModbusDevice) {
				defer wg.Done()
				// 如果设备连接错误次数超过 10， 并且最后链接时间小于当前5分钟，忽略设备
				if devitem.ConnErrTimes > 10 &&
					devitem.LastConnect.Add(time.Second*300).After(time.Now()) {
					return
				}
				switch devitem.ProtoType {
				case ConstProtoModbusRTU:
					glimit := limitPool.GetGoLimit(devitem.MbrtuAddr, 1)
					glimit.Add()
					defer glimit.Done()
					ReadModbusRtuRegData(devitem)
				case ConstProtoModbusTCP:
					glimit := limitPool.GetGoLimit(devitem.MbtcpAddr, 1)
					glimit.Add()
					defer glimit.Done()
					ReadModbusTcpRegData(devitem)
				}
			}(dev)
		}
		wg.Wait()
	}
}

// ReadModbusRtuRegData 读取 modbus rtu
func ReadModbusRtuRegData(dev models.ModbusDevice) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%+v", err)
		}
	}()
	var modbusregs []models.ModbusReg
	err := app.DB().Where("device_id = ?", dev.Id).Find(&modbusregs).Error
	if err != nil || modbusregs == nil || len(modbusregs) == 0 {
		// log.Errorf("读取设备[%s]寄存器列表失败或无寄存器... %s", dev.Name, err)
		time.Sleep(time.Millisecond * 1000)
		return
	}

	client, err := app.NewModbusRTUClient(dev.MbrtuAddr, dev.BaudRate, dev.PktDelay, dev.MbslaveId)
	if err != nil {
		log.Errorf("ModuBusRTU 设备[%s -> %s]连接失败 %s", dev.Name, dev.MbrtuAddr, err.Error())
		app.UpdateModbusDevConnStatus(dev.Id, err.Error())
		return
	}
	app.UpdateModbusDevConnStatus(dev.Id, "success")
	defer client.Close()

	for _, reg := range modbusregs {
		if reg.Status == "disabled" || common.InSlice(reg.DataType, []string{"Alarm", "Mark"}) {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if app.IsDebug() {
			log.Debugf("正在读取设备 [%s]-[%s] 数据 ", dev.Name, reg.Name)
		}

		switch reg.RegType {
		case "InputRegister":
			_read0304RegisterData(dev, reg, client, "InputRegister")
		case "HoldingRegister":
			_read0304RegisterData(dev, reg, client, "HoldingRegister")
		case "Coil":
			_readCoilData(dev, reg, client)
		case "DiscreteInput":
			_readDiscreteInputData(dev, reg, client)
		default:
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(reg.Intervals))
	}

}

// ReadModbusTcpRegData 读取 modbus tcp
func ReadModbusTcpRegData(dev models.ModbusDevice) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	var modbusregs []models.ModbusReg
	err := app.DB().Where("device_id = ?", dev.Id).Find(&modbusregs).Error
	if err != nil || modbusregs == nil || len(modbusregs) == 0 {
		// log.Errorf("读取设备[%s]寄存器列表失败或无寄存器... %s", dev.Name, err)
		time.Sleep(time.Millisecond * 1000)
		return
	}

	client, err := app.GetModbusTCPClient(dev.MbtcpAddr, dev.MbtcpPort, dev.MbslaveId)
	if err != nil {
		log.Errorf("ModuBusTCP 设备[%s -> %s:%d]连接失败 %s", dev.Name, dev.MbtcpAddr, dev.MbtcpPort, err.Error())
		app.UpdateModbusDevConnStatus(dev.Id, err.Error())
		return
	}
	app.UpdateModbusDevConnStatus(dev.Id, "success")
	defer client.Close()

	for _, reg := range modbusregs {
		if reg.Status == "disabled" || common.InSlice(reg.DataType, []string{"Alarm", "Mark"}) {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if log.IsDebug() {
			log.Debugf("正在读取设备 [%s]-[%s] 数据 ", dev.Name, reg.Name)
		}
		switch reg.RegType {
		case "InputRegister":
			_read0304RegisterData(dev, reg, client, "InputRegister")
		case "HoldingRegister":
			_read0304RegisterData(dev, reg, client, "HoldingRegister")
		case "Coil":
			_readCoilData(dev, reg, client)
		case "DiscreteInput":
			_readDiscreteInputData(dev, reg, client)
		default:
			continue
		}

		time.Sleep(time.Millisecond * time.Duration(reg.Intervals))
	}
}

// 读取modbus 寄存器数据并更新数据
func _readCoilData(dev models.ModbusDevice, reg models.ModbusReg, client modbus.Client) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	ret, err := client.ReadCoils(uint16(reg.StartAddr), 1)
	if err != nil {
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagFailure, err.Error())
		return
	}

	var val byte
	boolBuffer := bytes.NewReader(ret)
	err = binary.Read(boolBuffer, binary.BigEndian, &val)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, err.Error())
		return
	}
	if app.IsDebug() {
		log.Infof("读取到设备[%s]寄存器[%s]原始值 %d", dev.Name, reg.Name, val)
	}

	app.SaveModbusRegRtd(reg.Id, strconv.Itoa(int(val)), app.DataFlagSuccess, "success")
}

// 读取modbus 寄存器数据并更新数据
func _readDiscreteInputData(dev models.ModbusDevice, reg models.ModbusReg, client modbus.Client) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	ret, err := client.ReadDiscreteInputs(uint16(reg.StartAddr), 1)
	if err != nil {
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagFailure, err.Error())

		return
	}

	var val byte
	boolBuffer := bytes.NewReader(ret)
	err = binary.Read(boolBuffer, binary.BigEndian, &val)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, err.Error())
		return
	}

	if log.IsDebug() {
		log.Infof("读取到设备[%s]寄存器[%s]原始值 %d", dev.Name, reg.Name, val)
	}

	app.SaveModbusRegRtd(reg.Id, strconv.Itoa(int(val)), app.DataFlagSuccess, "success")

}

// 读取modbus 寄存器数据并更新数据
func _read0304RegisterData(dev models.ModbusDevice, reg models.ModbusReg, client modbus.Client, regtype string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	var regvar models.ModbusVar
	err := app.DB().Where("id = ?", reg.VarId).First(&regvar).Error
	if err != nil {
		log.Errorf("计算设备[%s]寄存器[%s]值失败, 未绑定变量", dev.Name, reg.Name)
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagStop, err.Error())
		return
	}

	var reglen uint16 = uint16(reg.DataLen / 2)
	var ret []byte
	switch regtype {
	case "HoldingRegister":
		ret, err = client.ReadHoldingRegisters(uint16(reg.StartAddr), reglen)
	case "InputRegister":
		ret, err = client.ReadInputRegisters(uint16(reg.StartAddr), reglen)
	default:
		err = errors.New("regtype error")
	}

	if err != nil {
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagFailure, err.Error())
		return
	}

	var val float64 = ERRORVALUE
	var dataval = make([]uint16, reglen)
	boolBuffer := bytes.NewReader(ret)
	err = binary.Read(boolBuffer, binary.BigEndian, &dataval)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, err.Error())
		return
	}

	// 如果寄存器个数是2,4, 则计算浮点数值
	switch reglen {
	case 1:
		val = float64(dataval[0])
	case 2:
		val = float64(modbus.GetFloat32Value(dataval[0], dataval[1], reg.GetByteOrder()))
	case 4:
		val = modbus.GetFloat64Value(dataval[0], dataval[1], dataval[2], dataval[3], reg.GetByteOrder())
	}

	if val == ERRORVALUE {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, "")
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, "ERRORVALUE")
		return
	}

	if app.IsDebug() {
		log.Debugf("读取到设备[%s]寄存器[%s]原始值 %v", dev.Name, reg.Name, dataval)
	}

	var result string

	// 模拟量计算
	if regvar.DxVal != "N/A" && regvar.DyVal != "N/A" {
		result, err = modbus.DoComputeDxyResult(val, regvar.DxVal, regvar.DyVal,
			regvar.MinVal, regvar.MaxVal, regvar.MaxAval, regvar.MinAval, int32(regvar.Decimals), regvar.Sign)
		if err != nil {
			log.Errorf("计算设备[%s]寄存器[%s]值失败, 模拟量计算错误, 请检查参数, %v", dev.Name, reg.Name, err.Error())
			app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, err.Error())
		}
		if app.IsDebug() {
			log.Infof("计算设备[%s]寄存器[%s]结果值 DoComputeDxyResult %s", dev.Name, reg.Name, result)
		}
	} else if reg.MaxSpval > 0 && reg.MaxSpval > reg.MinSpval {
		result, err = modbus.DoComputeResult(val, strconv.Itoa(reg.MinSpval), strconv.Itoa(reg.MaxSpval),
			regvar.MinVal, regvar.MaxVal, regvar.MaxAval, regvar.MinAval, int32(regvar.Decimals), regvar.Sign)
		if err != nil {
			log.Errorf("计算设备[%s]寄存器[%s]值失败, 模拟量计算错误, 请检查参数, %s", dev.Name, reg.Name, err.Error())
			app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagOver, err.Error())
		}
		if app.IsDebug() {
			log.Debugf("计算设备[%s]寄存器[%s]结果值 DoComputeResult %s", dev.Name, reg.Name, result)
		}
	}

	if result == "" {
		log.Errorf("计算设备[%s]寄存器[%s]值失败, 无结果, 请检查参数, %v", dev.Name, reg.Name, err)
		app.SaveModbusRegRtd(reg.Id, "N/A", app.DataFlagStop, "NoResult")
		return
	}

	app.SaveModbusRegRtd(reg.Id, result, app.DataFlagSuccess, "success")
}
