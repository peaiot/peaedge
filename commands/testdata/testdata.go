package main

import (
	"github.com/gocarina/gocsv"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/models"
)

func main() {
	app.Init(config.LoadConfig(""))
	app.Initdb()
	devs, err := readModbusDevice()
	common.Must(err)
	app.DB.Create(&devs)
	vars, err := readModbusVar()
	common.Must(err)
	app.DB.Create(&vars)
	regs, err := readModbusReg()
	common.Must(err)
	app.DB.Create(&regs)
}

func readModbusDevice() (data []*models.ModbusDevice, err error) {
	f, err := assets.TestData.Open("testdata/modbus_dev.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := gocsv.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func readModbusVar() (data []*models.ModbusVar, err error) {
	f, err := assets.TestData.Open("testdata/modbus_var.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := gocsv.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func readModbusReg() (data []*models.ModbusReg, err error) {
	f, err := assets.TestData.Open("testdata/modbus_reg.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := gocsv.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}
