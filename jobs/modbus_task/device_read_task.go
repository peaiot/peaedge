package modbus_task

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/golimit"
	"github.com/toughstruct/peaedge/common/goscript"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/models"
)

const (
	ConstProtoModbusRTU = "modbusrtu"
	ConstProtoModbusTCP = "modbustcp"
)

var sleep5s *time.Ticker
var sleep1s *time.Ticker
var sleep100ms *time.Ticker

func init() {
	sleep5s = time.NewTicker(time.Second * 5)
	sleep1s = time.NewTicker(time.Second * 1)
	sleep100ms = time.NewTicker(time.Millisecond * 100)
}

func StartModbusReadTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	<-sleep5s.C
	app.ResetModbusDevConnStatus()
	var limitPool = golimit.NewGoLimitPool()
	for {
		var devices = make([]models.ModbusDevice, 0)
		err := app.DB().Find(&devices).Error
		if err != nil {
			emsg := fmt.Sprintf("读取设备列表失败... %s", err.Error())
			log.Errorf(emsg)
			<-sleep5s.C
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
		<-sleep1s.C
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
		if reg.Status == "disabled" {
			<-sleep100ms.C
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
		<-sleep100ms.C
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
		<-sleep1s.C
		return
	}

	client, err := app.GetModbusTCPClient(dev.MbtcpAddr, dev.MbtcpPort, dev.MbslaveId)
	if err != nil {
		log.Errorf("ModuBusTCP 设备[%s -> %s:%d]连接失败 %s", dev.Name, dev.MbtcpAddr, dev.MbtcpPort, err.Error())
		app.UpdateModbusDevConnStatus(dev.Id, err.Error())
		return
	}
	app.UpdateModbusDevConnStatus(dev.Id, "success")

	for _, reg := range modbusregs {
		if reg.Status == "disabled" {
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

		time.Sleep(time.Millisecond * time.Duration(100))
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
		checkClientConn(err, client)
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagFailure, err.Error())
		return
	}

	var val byte
	boolBuffer := bytes.NewReader(ret)
	err = binary.Read(boolBuffer, binary.BigEndian, &val)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagOver, err.Error())
		return
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
		checkClientConn(err, client)
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagFailure, err.Error())
		return
	}

	var val byte
	boolBuffer := bytes.NewReader(ret)
	err = binary.Read(boolBuffer, binary.BigEndian, &val)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagOver, err.Error())
		return
	}

	app.SaveModbusRegRtd(reg.Id, strconv.Itoa(int(val)), app.DataFlagSuccess, "success")

}

// 读取 modbus 寄存器数据并更新数据
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
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagStop, err.Error())
		return
	}

	var reglen = uint16(reg.DataLen)
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
		checkClientConn(err, client)
		log.Errorf("读取设备[%s]寄存器[%s]失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagFailure, err.Error())
		return
	}

	var result = common.NA
	var dataval = make([]uint16, reglen)
	err = binary.Read(bytes.NewReader(ret), binary.BigEndian, &dataval)
	if err != nil {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, err.Error())
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagOver, err.Error())
		return
	}

	// 如果寄存器个数是2,4, 则计算浮点数值
	switch regvar.DataType {
	case "int16":
		result = cast.ToString(int16(ret[0])<<8 | int16(ret[1]))
	case "uint16":
		result = cast.ToString(uint16(ret[0])<<8 | uint16(ret[1]))
	case "float":
		result = cast.ToString(modbus.GetFloat32Value(dataval[0], dataval[1], regvar.GetByteOrder()))
	case "float64":
		result = cast.ToString(modbus.GetFloat64Value(dataval[0], dataval[1], dataval[2], dataval[3], regvar.GetByteOrder()))
	}

	if result == common.NA {
		log.Errorf("解析设备[%s]寄存器[%s]值失败 %s", dev.Name, reg.Name, "")
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagOver, "ERRORVALUE")
		return
	}

	// 数值计算
	var script models.DataScript
	err = app.DB().Where("id = ?", regvar.ScriptId).First(&script).Error
	if err == nil {
		vret, err := goscript.RunFunc(script.Content, script.FuncName, cast.ToFloat64(result))
		if err != nil {
			app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagOver, err.Error())
			log.Errorf("计算设备[%s]寄存器[%s]值失败, 请检查参数, %v", dev.Name, reg.Name, err.Error())
			return
		}
		result = cast.ToString(vret)
	}

	if result == "" || result == common.NA {
		log.Errorf("计算设备[%s]寄存器[%s]值失败, 无结果, 请检查参数, %v", dev.Name, reg.Name, err)
		app.SaveModbusRegRtd(reg.Id, common.NA, app.DataFlagStop, "NoResult")
		return
	}

	app.SaveModbusRegRtd(reg.Id, result, app.DataFlagSuccess, "success")
}

func checkClientConn(err error, client modbus.Client) {
	if errors.Is(err, net.ErrClosed) ||
		errors.Is(err, net.ErrWriteToConnected) ||
		errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, io.EOF) {
		_ = client.ReConnect()
	}
}
