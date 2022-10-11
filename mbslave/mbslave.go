package mbslave

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/goburrow/serial"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/common/mbserver"
	"github.com/toughstruct/peaedge/models"
)

var slave *ModbusSlave

type ModbusSlave struct {
	root *mbserver.Server
}

func NewModbusSlave() *ModbusSlave {
	return &ModbusSlave{}
}

func Listen() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	slave = NewModbusSlave()
	slave.root = mbserver.NewServer()
	tcpAddr := app.Config().Modbus.TcpAddr
	if tcpAddr != "" {
		err := slave.root.ListenTCP(tcpAddr)
		if err != nil {
			log.Error("modbus slave listen tcp %s error: %v", tcpAddr, err)
		}
	}
	rtuAddr := app.Config().Modbus.RtuAddr
	if !common.InSlice(rtuAddr, []string{"none", "null", "", "/dev/null"}) {
		err := slave.root.ListenRTU(&serial.Config{
			Address:  "/dev/ttyACM0",
			BaudRate: app.Config().Modbus.BaudRate,
			DataBits: app.Config().Modbus.DataBits,
			StopBits: app.Config().Modbus.StopBits,
			Parity:   app.Config().Modbus.Parity,
			Timeout:  time.Millisecond * time.Duration(app.Config().Modbus.Timeout),
			RS485: serial.RS485Config{
				Enabled:            true,
				DelayRtsBeforeSend: 2 * time.Millisecond,
				DelayRtsAfterSend:  3 * time.Millisecond,
				RtsHighDuringSend:  false,
				RtsHighAfterSend:   false,
				RxDuringTx:         false,
			}})
		if err != nil {
			log.Error("modbus slave listen rtu %s error: %v", rtuAddr, err)
		}
	}

	err := ReloadRegisterData()
	if err != nil {
		log.Error("modbus slave reload register data error: %v", err)
	}

	return nil
}

func ReloadRegisterData() (err error) {
	var datas []models.ModbusSlaveReg
	err = app.DB().Find(&datas).Error
	if err != nil {
		return err
	}
	for _, data := range datas {
		hbyte, err := hex.DecodeString(data.Value)
		if err != nil {
			log.Errorf("modbus slave reload register %s(%d) data (%s) error: %v",
				data.Name, data.Register, data.Value, err)
			continue
		}
		switch data.RegType {
		case "Coil":
			SetCoilData(uint16(data.Register), hbyte...)
		case "HoldingRegister":
			SetHoldingRegisterData(uint16(data.Register), mbserver.BytesToUint16(hbyte)...)
		case "InputRegister":
			SetInputRegisterData(uint16(data.Register), mbserver.BytesToUint16(hbyte)...)
		case "DiscreteInput":
			SetDiscreteInputData(uint16(data.Register), hbyte...)
		}
	}
	return nil
}

func SetHoldingRegisterData(register uint16, value ...uint16) {
	if register < 0 || register > 65535 {
		return
	}
	copy(slave.root.HoldingRegisters[register:], value)
}

func SetCoilData(register uint16, value ...byte) {
	if register < 0 || register > 65535 {
		return
	}
	copy(slave.root.Coils[register:], value)
}

func SetInputRegisterData(register uint16, value ...uint16) {
	if register < 0 || register > 65535 {
		return
	}
	copy(slave.root.InputRegisters[register:], value)
}

func SetDiscreteInputData(register uint16, value ...byte) {
	if register < 0 || register > 65535 {
		return
	}
	copy(slave.root.DiscreteInputs[register:], value)
}
