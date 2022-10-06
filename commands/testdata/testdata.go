package main

import (
	"time"

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
	app.DB().Create(&devs)
	vars, err := readModbusVar()
	common.Must(err)
	app.DB().Create(&vars)
	regs, err := readModbusReg()
	common.Must(err)
	app.DB().Create(&regs)
	oprs, err := readSysopr()
	common.Must(err)
	for _, opr := range oprs {
		opr.Password = common.Sha256Hash(opr.Password)
	}
	app.DB().Create(&oprs)

	ds := models.DataScript{
		ID:        "bm3wnhieckryy",
		Name:      "modbus 值计算",
		FuncName:  "HandlerModbusData",
		Content:   assets.ModbusDataFuncs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	app.DB().Create(&ds)
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

func readSysopr() (data []*models.SysOpr, err error) {
	f, err := assets.TestData.Open("testdata/sys_opr.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := gocsv.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}
