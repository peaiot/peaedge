package app

import (
	"time"

	"github.com/toughstruct/peaedge/common/log"
)

// SaveModbusRegRtd 保存寄存器数据
//goland:noinspection SqlDialectInspection
func SaveModbusRegRtd(regId string, result string, flag, errstr string) {
	var err error
	if errstr == "" || errstr == "success" {
		err = gormDB.Exec("update modbus_reg set rtd = ?,flag = ?, err_times = 0, last_error= 'success'  where id = ?",
			result, flag, regId).Error
	} else {
		err = gormDB.Exec("update modbus_reg set rtd = ?, flag = ?, err_times = IFNULL(err_times,0) +1, last_error= ?  where id = ?",
			result, flag, errstr, regId).Error
	}
	if err != nil {
		log.Errorf("SaveModbusRegRtd error: %s", err)
	}
}

func UpdateModbusDevConnStatus(devId string, errstr string) {
	var err error
	if errstr == "" || errstr == "success" {
		err = gormDB.Exec("update modbus_device set last_connect = ?, conn_err_times = 0, last_conn_error = 'success'  where id = ?",
			time.Now(), devId).Error
	} else {
		err = gormDB.Exec("update modbus_device set last_connect = ?, conn_err_times = IFNULL(conn_err_times, 0) +1, last_conn_error = ?  where id = ?",
			time.Now(), errstr, devId).Error
	}
	if err != nil {
		log.Errorf("UpdateModbusDevConnStatus error: %s", err)
	}
}

func ResetModbusDevConnStatus() {
	var err error
	err = gormDB.Exec("update modbus_device set last_connect = ?, conn_err_times = 0, last_conn_error = 'success'", time.Now()).Error
	if err != nil {
		log.Errorf("ResetModbusDevConnStatus error: %s", err)
	}
}
