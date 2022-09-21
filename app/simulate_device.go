package app

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"regexp"
	"strings"

	"github.com/toughstruct/peaedge/models"
)

var testdataReg = regexp.MustCompile(`\[testdata=(.*)\]`)

type SimulateDevice struct {
	Proto   string
	Addr    string
	TcpPort int
	SlaveId int
}

func NewSimulateDevice(proto string, addr string, tcpPort int, slaveId int) *SimulateDevice {
	return &SimulateDevice{Proto: proto, Addr: addr, TcpPort: tcpPort, SlaveId: slaveId}
}

func (s SimulateDevice) getModbusReg(startAddr uint16) (re *models.ModbusReg, err error) {
	var dev models.ModbusDevice
	switch s.Proto {
	case "modbusrtu":
		err = DB.Where("proto_type = ? and mbrtu_addr = ? and mbslave_id = ?", s.Proto, s.Addr, s.SlaveId).First(&dev).Error
	case "modbustcp":
		err = DB.Where("proto_type = ? and mbtcp_addr = ? and mbtcp_port = ? and mbslave_id = ?", s.Proto, s.Addr, s.SlaveId).First(&dev).Error
	}
	if err != nil {
		return nil, err
	}
	var reg models.ModbusReg
	err = DB.Where("device_id = ? and start_addr = ?", dev.Id, startAddr).First(&reg).Error
	return &reg, err
}

var errtestdata = errors.New("测试数据定义错误")

func (s SimulateDevice) ReadTestData(address uint16) (results []byte, err error) {
	reg, err := s.getModbusReg(address)
	if err != nil {
		return nil, err
	}
	subv := testdataReg.FindStringSubmatch(reg.Remark)
	if subv != nil {
		values := strings.Split(subv[1], ",")
		val := values[rand.Intn(len(values))]
		rdata, err := hex.DecodeString(val)
		if err != nil {
			return nil, err
		}
		return rdata, nil
	}
	return nil, errtestdata
}

func (s SimulateDevice) ReadCoils(address, quantity uint16) (results []byte, err error) {
	return s.ReadTestData(address)
}

func (s SimulateDevice) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	return s.ReadTestData(address)
}

func (s SimulateDevice) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	return s.ReadTestData(address)
}

func (s SimulateDevice) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return s.ReadTestData(address)
}

func (s SimulateDevice) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) ReadFIFOQueue(address uint16) (results []byte, err error) {
	return nil, nil
}

func (s SimulateDevice) Connect() error {
	return nil
}

func (s SimulateDevice) Close() error {
	return nil
}

func (s SimulateDevice) ReConnect() error {
	return nil
}
