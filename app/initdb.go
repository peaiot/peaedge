package app

import (
	"time"

	"github.com/gocarina/gocsv"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/models"
)

func Initdb() {
	_ = gormDB.Migrator().DropTable(models.Tables...)
	_ = gormDB.Migrator().AutoMigrate(models.Tables...)
	initSysData()
}

func initSysData() {
	gormDB.Create(&models.SysOpr{
		ID:        "10000",
		Realname:  "系统管理员",
		Mobile:    "000000",
		Email:     "master@peaiot.net",
		Username:  "admin",
		Password:  common.Sha256Hash("peaedge"),
		Level:     "super",
		Status:    "enabled",
		Remark:    "系统管理员",
		LastLogin: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	ds := []models.DataScript{
		{
			ID:        "1000001",
			Name:      "modbus 值计算",
			FuncName:  "HandlerModbusData",
			Content:   assets.HandlerModbusData,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "1000002",
			Name:      "数据流脚本",
			FuncName:  "HandlerDataStream",
			Content:   assets.HandlerDataStream,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	gormDB.Create(&ds)
}

func InitTestData() {
	devs, err := readTestModbusDevice()
	common.Must(err)
	gormDB.Save(&devs)

	vars, err := readTestModbusVar()
	common.Must(err)
	gormDB.Save(&vars)

	regs, err := readTestModbusReg()
	common.Must(err)
	gormDB.Save(&regs)

	msregs, err := readTestModbusSlaveReg()
	common.Must(err)
	gormDB.Save(&msregs)
}

func readTestModbusDevice() (data []*models.ModbusDevice, err error) {
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

func readTestModbusVar() (data []*models.ModbusVar, err error) {
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

func readTestModbusReg() (data []*models.ModbusReg, err error) {
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

func readTestModbusSlaveReg() (data []*models.ModbusSlaveReg, err error) {
	f, err := assets.TestData.Open("testdata/modbus_slave_reg.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := gocsv.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}
