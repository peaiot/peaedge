package app

import (
	"fmt"
	slog "log"
	"os"
	"time"

	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/common/modbus"
)

// NewModbusRTUClient 创建新的 modbusrtu 客户端
func NewModbusRTUClient(rtuaddr string, baudRate int, pktDelay int, slaveId int) (modbus.Client, error) {
	// 模拟设备
	if os.Getenv("PEAEDGE_VMDEV") != "" {
		return NewSimulateDevice("modbusrtu", rtuaddr, 0, slaveId), nil
	}
	handler := modbus.NewRTUClientHandler(rtuaddr)
	handler.SlaveId = byte(slaveId)
	handler.BaudRate = baudRate
	handler.MsgDelay = time.Duration(pktDelay)
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.IdleTimeout = time.Second * 5
	handler.Timeout = 1000 * time.Millisecond
	if IsDebug() {
		handler.Logger = slog.New(log.Stdlog{}, "Modbus-RTU", 0)
	}
	err := handler.Connect()
	if err != nil {
		return nil, err
	}
	return modbus.NewClient(handler), nil
}

// GetModbusTCPClient 获取TCP 客户端
func GetModbusTCPClient(devaddr string, port int, slaveid int) (modbus.Client, error) {
	// 模拟设备
	if os.Getenv("PEAEDGE_VMDEV") != "" {
		return NewSimulateDevice("modbustcp", devaddr, port, slaveid), nil
	}
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", devaddr, port))
	handler.SlaveId = byte(slaveid)
	handler.Timeout = 1000 * time.Millisecond
	handler.IdleTimeout = time.Second * 10
	if IsDebug() {
		handler.Logger = slog.New(log.Stdlog{}, "Modbus-TCP", 0)
	}
	err := handler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(handler)
	return client, nil
}
