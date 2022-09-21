package models

import (
	"time"

	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/modbus"
)

// SysConfig 系统配置
type SysConfig struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Remark    string `json:"remark"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ModbusDevice Modbus 设备
type ModbusDevice struct {
	Id            string    `json:"id" csv:"id"`
	Name          string    `json:"name" csv:"name"`
	MN            string    `json:"mn" csv:"mn"`
	ProtoType     string    `json:"proto_type" csv:"proto_type"`
	MbrtuAddr     string    `json:"mbrtu_addr" csv:"mbrtu_addr"`
	MbtcpAddr     string    `json:"mbtcp_addr" csv:"mbtcp_addr"`
	MbtcpPort     int       `json:"mbtcp_port" csv:"mbtcp_port"`
	MbslaveId     int       `json:"mbslave_id" csv:"mbslave_id"`
	BaudRate      int       `json:"baud_rate" csv:"baud_rate"`
	PktDelay      int       `json:"pkt_delay" csv:"pkt_delay"`
	Remark        string    `json:"remark" csv:"remark"`
	ConnErrTimes  int       `json:"conn_err_times"`
	LastConnError string    `json:"last_conn_error"`
	LastConnect   time.Time `json:"last_connect"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ModbusReg Modbus寄存器
type ModbusReg struct {
	Id         string    `json:"id" csv:"id"`
	DeviceId   string    `json:"device_id" csv:"device_id"`
	Name       string    `json:"name" csv:"name"`
	DataType   string    `json:"data_type" csv:"data_type"`
	RegType    string    `json:"reg_type" csv:"reg_type"`
	StartAddr  int       `json:"start_addr" csv:"start_addr"`
	ByteAddr   int       `json:"byte_addr" csv:"byte_addr"`
	BitAddr    int       `json:"bit_addr" csv:"bit_addr"`
	DataLen    int       `json:"data_len" csv:"data_len"`
	Intervals  int       `json:"intervals" csv:"intervals"`
	Decimals   int       `json:"decimals" csv:"decimals"`
	ByteOrder  string    `json:"byte_order" csv:"byte_order"`
	AccessType string    `json:"access_type" csv:"access_type"`
	MinSpval   int       `json:"min_spval" csv:"min_spval"`
	MaxSpval   int       `json:"max_spval" csv:"max_spval"`
	VarId      string    `json:"var_id" csv:"var_id"`
	Rtd        string    `json:"rtd" csv:"rtd"`
	LastUpdate time.Time `json:"last_update" csv:"last_update"`
	ErrTimes   int       `json:"err_times"`
	LastError  string    `json:"last_error"`
	Status     string    `json:"status" csv:"status"`
	Remark     string    `json:"remark" csv:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ModbusVar 变量定义
type ModbusVar struct {
	Id             string    `json:"id" csv:"id"`
	Name           string    `json:"name" csv:"name"`
	DataType       string    `json:"data_type" csv:"data_type"`
	Unit           string    `json:"unit" csv:"unit"`
	InitVal        string    `json:"init_val" csv:"init_val"`
	MinVal         string    `json:"min_val" csv:"min_val"`
	MaxVal         string    `json:"max_val" csv:"max_val"`
	MinAval        string    `json:"min_aval" csv:"min_aval"`
	MaxAval        string    `json:"max_aval" csv:"max_aval"`
	DxVal          string    `json:"dx_val" csv:"dx_val"`
	DyVal          string    `json:"dy_val" csv:"dy_val"`
	SaveDelay      int       `json:"save_delay" csv:"save_delay"`
	Decimals       int       `json:"decimals" csv:"decimals"`
	Sign           int       `json:"sign" csv:"sign"`
	Jscript        string    `json:"jscript" csv:"jscript"`
	Hj212Attr      string    `json:"hj212_attr" csv:"hj212_attr"`
	Hj212RtdStatus string    `json:"hj212_rtd_status" csv:"hj212_rtd_status"`
	Remark         string    `json:"remark" csv:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Hj212Queue 212消息队列
type Hj212Queue struct {
	Id         string    `json:"id"`
	Servid     string    `json:"servid"`
	Cn         string    `json:"cn"`
	CpMessage  string    `json:"cp_message"`
	SendStatus string    `json:"send_status"`
	SendTimes  int       `json:"send_times"`
	CreateTime time.Time `json:"create_time"`
	LastSend   time.Time `json:"last_send"`
}

// Hj212Serv 212服务器
type Hj212Serv struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Server    string    `json:"server"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IotDevice 通用物联网设备
type IotDevice struct {
	ID        string    `json:"id" form:"id"`
	MN        string    `gorm:"uniqueIndex" json:"mn" form:"mn"`
	Name      string    `json:"name" form:"name"`
	Remark    string    `json:"remark" form:"remark"`
	ProtoType string    `json:"proto_type" form:"device_type"` // 协议类型 zigbee ble lora
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeviceRtdData 设备实时数据
type DeviceRtdData struct {
	ID        string    `json:"id"`
	MN        string    `json:"mn"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

var Tables = []interface{}{
	&SysConfig{},
	// hj212
	&Hj212Serv{},
	&Hj212Queue{},
	// modbus
	&ModbusDevice{},
	&ModbusReg{},
	&ModbusVar{},
	// iot device
	&IotDevice{},
	&DeviceRtdData{},
}

func (r *ModbusReg) GetByteOrder() string {
	if !common.InSlice(r.ByteOrder, []string{modbus.BigEndian, modbus.BigEndianSwap, modbus.LittleEndian, modbus.LittleEndianSwap}) {
		return modbus.BigEndian
	}
	return r.ByteOrder
}
