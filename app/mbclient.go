package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/log"
)

// NewModbusRTUClient 创建新的 modbusrtu 客户端
func NewModbusRTUClient(rtuaddr string, baudRate int, pktDelay int, slaveId int) (modbus.Client, error) {
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
		handler.Logger = log.Modbus
	}
	err := handler.Connect()
	if err != nil {
		return nil, err
	}
	return modbus.NewClient(handler), nil
}

var tcpMbPoolLock sync.Mutex
var tcpMbPool map[string]modbus.Client

func init() {
	tcpMbPoolLock = sync.Mutex{}
	tcpMbPool = make(map[string]modbus.Client)
}

// GetModbusTCPClient 获取TCP 客户端
func GetModbusTCPClient(devaddr string, port int, slaveid int) (modbus.Client, error) {
	tcpMbPoolLock.Lock()
	defer tcpMbPoolLock.Unlock()
	key := fmt.Sprintf("tcp://%s:%d/%d", devaddr, port, slaveid)
	if c, ok := tcpMbPool[key]; ok {
		return c, nil
	}
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", devaddr, port))
	handler.SlaveId = byte(slaveid)
	handler.Timeout = 1000 * time.Millisecond
	handler.IdleTimeout = time.Hour * 24
	if IsDebug() {
		handler.Logger = log.Modbus
	}
	err := handler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(handler)
	tcpMbPool[key] = client
	return client, nil
}
