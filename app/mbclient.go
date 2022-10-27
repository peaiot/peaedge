package app

import (
	"fmt"
	"time"

	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/log"
)

// // NewModbusRTUClient 创建新的 modbusrtu 客户端
// func NewModbusRTUClient(rtuaddr string, baudRate int, pktDelay int, slaveId int) (modbus.Client, error) {
// 	handler := modbus.NewRTUClientHandler(rtuaddr)
// 	handler.SlaveId = byte(slaveId)
// 	handler.BaudRate = baudRate
// 	handler.MsgDelay = time.Duration(pktDelay)
// 	handler.DataBits = 8
// 	handler.Parity = "N"
// 	handler.StopBits = 1
// 	handler.IdleTimeout = time.Second * 5
// 	handler.Timeout = 1000 * time.Millisecond
// 	if IsDebug() {
// 		handler.Logger = log.Modbus
// 	}
// 	err := handler.Connect()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return modbus.NewClient(handler), nil
// }

// GetModbusRTUClient 获取 RTU 客户端
func GetModbusRTUClient(address string, baudRate int, pktDelay int, slaveId int) (modbus.Client, error) {
	transporter := mbRtuPool.Get(address)
	if IsDebug() {
		transporter.Logger = log.Modbus
	}
	transporter.BaudRate = baudRate
	transporter.MsgDelay = time.Duration(pktDelay)
	transporter.DataBits = 8
	transporter.Parity = "N"
	transporter.StopBits = 1
	transporter.IdleTimeout = time.Hour * 24
	transporter.Timeout = 1000 * time.Millisecond
	err := transporter.Connect()
	if err != nil {
		return nil, err
	}
	packager := modbus.NewRTUClientPackager(slaveId)
	client := modbus.NewClient2(packager, transporter)
	return client, nil
}

func ReleaseModbusRTUClient(address string, c modbus.Client) {
	mbRtuPool.Put(address, c.GetTransporter().(*modbus.RTUClientTransporter))
}

// GetModbusTCPClient 获取TCP 客户端
func GetModbusTCPClient(devaddr string, port int, slaveid int) (modbus.Client, error) {
	address := fmt.Sprintf("%s:%d", devaddr, port)
	transporter := mbTcpPool.Get(address)
	if IsDebug() {
		transporter.Logger = log.Modbus
	}
	err := transporter.Connect()
	if err != nil {
		return nil, err
	}
	packager := modbus.NewTCPClientPackager(slaveid)
	client := modbus.NewClient2(packager, transporter)
	return client, nil
}

func ReleaseModbusTCPClient(address string, c modbus.Client) {
	mbTcpPool.Put(address, c.GetTransporter().(*modbus.TCPClientTransporter))
}
