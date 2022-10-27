package models

import (
	"encoding/json"
	"fmt"

	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/common/timeutil"
)

func (d ModbusReg) MarshalJSON() ([]byte, error) {
	type Alias ModbusReg
	return json.Marshal(&struct {
		Alias
		LastUpdate string `json:"last_update"`
	}{
		Alias:      (Alias)(d),
		LastUpdate: timeutil.FmtDatetimeString(d.LastUpdate),
	})
}

func (d Hj212Queue) MarshalJSON() ([]byte, error) {
	type Alias Hj212Queue
	return json.Marshal(&struct {
		Alias
		CreateTime string `json:"create_time"`
		LastSend   string `json:"last_send"`
	}{
		Alias:      (Alias)(d),
		CreateTime: timeutil.FmtDatetimeString(d.CreateTime),
		LastSend:   timeutil.FmtDatetimeString(d.LastSend),
	})
}

func (r ModbusSlaveReg) Sign() string {
	return fmt.Sprintf("%s:%d", r.RegType, r.Register)
}

func (r *ModbusVar) GetByteOrder() string {
	if !common.InSlice(r.ByteOrder, []string{modbus.BigEndian, modbus.BigEndianSwap, modbus.LittleEndian, modbus.LittleEndianSwap}) {
		return modbus.BigEndian
	}
	return r.ByteOrder
}
