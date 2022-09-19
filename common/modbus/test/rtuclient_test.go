// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/dpapathanasiou/go-modbus"
	"github.com/goburrow/modbus"
)

const (
	rtuDevice = "/dev/pts/6"
)

func TestRTUClient(t *testing.T) {
	// Diagslave does not support broadcast id.
	handler := modbus.NewRTUClientHandler(rtuDevice)
	handler.SlaveId = 17
	ClientTestAll(t, modbus.NewClient(handler))
}

func TestRTUClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewRTUClientHandler(rtuDevice)
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(15, 2)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	results, err = client.ReadWriteMultipleRegisters(0, 2, 2, 2, []byte{1, 2, 3, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}

func TestRTUClientRead(t *testing.T) {
	handler := modbus.NewRTUClientHandler("/dev/cu.usbserial-14610")
	handler.BaudRate = 9600
	handler.Parity = "N"
	handler.DataBits = 8
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.IdleTimeout = 10 * time.Second
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(uint16(0), uint16(3))
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}

func TestXXX(t *testing.T) {
	ctx, _ := modbusclient.ConnectRTU("/dev/cu.usbserial-14610", 4800)
	readResult, readErr := modbusclient.RTURead(ctx, byte(1), modbusclient.FUNCTION_READ_HOLDING_REGISTERS, uint16(1), uint16(3), 100, true)
	if readErr != nil {
		log.Println(readErr)
	}
	log.Println(readResult)
}
