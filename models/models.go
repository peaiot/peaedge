package models

import (
	"time"
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
	Id        string `json:"id"`
	Name      string `json:"name"`
	MN        string `gorm:"uniqueIndex" json:"mn"`
	ProtoType string `json:"proto_type"`
	MbrtuAddr string `json:"mbrtu_addr"`
	MbtcpAddr string `json:"mbtcp_addr"`
	MbtcpPort int    `json:"mbtcp_port"`
	MbslaveId int    `json:"mbslave_id"`
	BaudRate  int    `json:"baud_rate"`
	PktDelay  int    `json:"pkt_delay"`
	Remark    string `json:"remark"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ModbusReg Modbus寄存器
type ModbusReg struct {
	Id         string    `json:"id"`
	DeviceId   string    `json:"device_id"`
	Name       string    `json:"name"`
	DataType   string    `json:"data_type"`
	RegType    string    `json:"reg_type"`
	StartAddr  int       `json:"start_addr"`
	ByteAddr   int       `json:"byte_addr"`
	BitAddr    int       `json:"bit_addr"`
	DataLen    int       `json:"data_len"`
	Intervals  int       `json:"intervals"`
	Decimals   int       `json:"decimals"`
	ByteOrder  string    `json:"byte_order"`
	AccessType string    `json:"access_type"`
	MinSpval   int       `json:"min_spval"`
	MaxSpval   int       `json:"max_spval"`
	VarId      string    `json:"var_id"`
	Rtd        string    `json:"rtd"`
	LastUpdate time.Time `json:"last_update"`
	Status     string    `json:"status"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ModbusVar 变量定义
type ModbusVar struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	DataType      string `json:"data_type"`
	Unit          string `json:"unit"`
	InitVal       string `json:"init_val"`
	MinVal        string `json:"min_val"`
	MaxVal        string `json:"max_val"`
	MinAval       string `json:"min_aval"`
	MaxAval       string `json:"max_aval"`
	DxVal         string `json:"dx_val"`
	DyVal         string `json:"dy_val"`
	SaveDelay     int    `json:"save_delay"`
	Decimals      int    `json:"decimals"`
	Sign          int    `json:"sign"`
	Jscript       string `json:"jscript"`
	H212Attr      string `json:"h212_attr"`
	H212RtdStatus string `json:"h212_rtd_status"`
	Remark        string `json:"remark"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
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
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Server    string `json:"server"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
	ID        string    `gorm:"primaryKey" json:"id"`
	MN        string    `json:"mn"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
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
