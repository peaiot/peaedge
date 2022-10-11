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

type SysOpr struct {
	ID        int64     `json:"id,string" form:"id"  csv:"id"`
	Realname  string    `json:"realname" form:"realname" csv:"realname"`
	Mobile    string    `json:"mobile" form:"mobile" csv:"mobile"`
	Email     string    `json:"email" form:"email" csv:"email"`
	Username  string    `json:"username" form:"username" csv:"username"`
	Password  string    `json:"password" form:"password" csv:"password"`
	Level     string    `json:"level" form:"level" csv:"level"`
	Status    string    `json:"status" form:"status" csv:"status"`
	Remark    string    `json:"remark" form:"remark" csv:"remark"`
	LastLogin time.Time `json:"last_login" form:"last_login" csv:"last_login"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt time.Time `json:"updated_at" csv:"updated_at"`
}

type DataScript struct {
	ID        string    `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	FuncName  string    `json:"func_name" form:"func_name"`
	Content   string    `json:"content" form:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ModbusDevice Modbus 设备
type ModbusDevice struct {
	Id            string    `json:"id" csv:"id" form:"id"`
	Name          string    `json:"name" csv:"name" form:"name"`
	MN            string    `json:"mn" csv:"mn" form:"mn"`
	ProtoType     string    `json:"proto_type" csv:"proto_type" form:"proto_type"`
	MbrtuAddr     string    `json:"mbrtu_addr" csv:"mbrtu_addr" form:"mbrtu_addr"`
	MbtcpAddr     string    `json:"mbtcp_addr" csv:"mbtcp_addr" form:"mbtcp_addr"`
	MbtcpPort     int       `json:"mbtcp_port" csv:"mbtcp_port" form:"mbtcp_port"`
	MbslaveId     int       `json:"mbslave_id" csv:"mbslave_id" form:"mbslave_id"`
	BaudRate      int       `json:"baud_rate" csv:"baud_rate" form:"baud_rate"`
	PktDelay      int       `json:"pkt_delay" csv:"pkt_delay" form:"pkt_delay"`
	Remark        string    `json:"remark" csv:"remark" form:"remark"`
	ConnErrTimes  int       `json:"conn_err_times" form:"conn_err_times"`
	LastConnError string    `json:"last_conn_error" form:"last_conn_error"`
	LastConnect   time.Time `json:"last_connect" form:"last_connect"`
	CreatedAt     time.Time `json:"created_at" `
	UpdatedAt     time.Time `json:"updated_at"`
}

// ModbusReg Modbus寄存器
type ModbusReg struct {
	Id         string    `json:"id" csv:"id" form:"id"`
	DeviceId   string    `json:"device_id" csv:"device_id" form:"device_id"`
	Name       string    `json:"name" csv:"name" form:"name"`
	RegType    string    `json:"reg_type" csv:"reg_type" form:"reg_type"`
	StartAddr  int       `json:"start_addr" csv:"start_addr" form:"start_addr"`
	DataLen    int       `json:"data_len" csv:"data_len" form:"data_len"`
	AccessType string    `json:"access_type" csv:"access_type" form:"access_type"`
	Rtd        string    `json:"rtd" csv:"rtd" form:"rtd"`
	Flag       string    `json:"flag" csv:"flag" form:"flag"`
	VarId      string    `json:"var_id" csv:"var_id" form:"var_id"`
	LastUpdate time.Time `json:"last_update" csv:"last_update"`
	ErrTimes   int       `json:"err_times" form:"err_times"`
	LastError  string    `json:"last_error" form:"last_error"`
	Status     string    `json:"status" csv:"status" form:"status"`
	Remark     string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ModbusVar 变量定义
type ModbusVar struct {
	Id            string    `json:"id" csv:"id" form:"id"`
	Name          string    `json:"name" csv:"name" form:"name"`
	DataType      string    `json:"data_type" csv:"data_type" form:"data_type"`
	Unit          string    `json:"unit" csv:"unit" form:"unit"`
	ByteOrder     string    `json:"byte_order" csv:"byte_order" form:"byte_order"`
	ScriptId      string    `json:"script_id" csv:"script_id" form:"script_id"`
	DataFactor    string    `json:"data_factor" csv:"data_factor" form:"data_factor"`
	ChannelStatus string    `json:"channel_status" csv:"channel_status" form:"channel_status"`
	Remark        string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" `
}

type ModbusVarTpl struct {
	Id            string `json:"id" csv:"id" from:"id"`
	Name          string `json:"name" csv:"name" from:"name"`
	DataType      string `json:"data_type" csv:"data_type" from:"data_type"`
	Unit          string `json:"unit" csv:"unit" from:"unit"`
	ByteOrder     string `json:"byte_order" csv:"byte_order" from:"byte_order"`
	ScriptId      string `json:"script_id" csv:"script_id" from:"script_id"`
	DataFactor    string `json:"data_factor" csv:"data_factor" form:"data_factor"`
	ChannelStatus string `json:"channel_status" csv:"channel_status" form:"channel_status"`
	Remark        string `json:"remark" csv:"remark" from:"remark"`
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
	MN        string    `json:"mn"`         // 机器 MN
	Name      string    `json:"name"`       // 数据属性
	Factor    string    `json:"factor"`     // 环境因子
	Value     string    `json:"value"`      // 实时值
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// ModbusSlaveReg modbus server 寄存器，系统启动时加载到内存
type ModbusSlaveReg struct {
	ID       string `json:"id" csv:"id" form:"id"`
	Name     string `json:"name" csv:"name" form:"name"`
	Register int    `json:"register" csv:"register" form:"register"`
	RegType  string `json:"reg_type" csv:"reg_type" form:"reg_type"`
	Length   int    `json:"length" csv:"length" form:"length"`
	Value    string `json:"value" csv:"value" form:"value"`
	Remark   string `json:"remark" csv:"remark" form:"remark"`
}

var Tables = []interface{}{
	&SysConfig{},
	&SysOpr{},
	&DataScript{},
	// hj212
	&Hj212Serv{},
	&Hj212Queue{},
	// modbus
	&ModbusDevice{},
	&ModbusReg{},
	&ModbusVar{},
	&ModbusVarTpl{},
	// iot device
	&IotDevice{},
	&DeviceRtdData{},
	&ModbusSlaveReg{},
}

func (r *ModbusVar) GetByteOrder() string {
	if !common.InSlice(r.ByteOrder, []string{modbus.BigEndian, modbus.BigEndianSwap, modbus.LittleEndian, modbus.LittleEndianSwap}) {
		return modbus.BigEndian
	}
	return r.ByteOrder
}
