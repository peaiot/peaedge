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

type SysOpr struct {
	ID        string    `json:"id" form:"id"  csv:"id"`
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
	Oid           string    `json:"oid" csv:"oid" form:"oid"`
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
	Oid        string    `json:"oid" csv:"oid" form:"oid"`
	RegType    string    `json:"reg_type" csv:"reg_type" form:"reg_type"`
	StartAddr  int       `json:"start_addr" csv:"start_addr" form:"start_addr"`
	AccessType string    `json:"access_type" csv:"access_type" form:"access_type"`
	Rtd        string    `json:"rtd" csv:"rtd" form:"rtd"`
	Flag       string    `json:"flag" csv:"flag" form:"flag"`
	VarId      string    `json:"var_id" csv:"var_id" form:"var_id"`
	DataType   string    `json:"data_type" csv:"data_type" form:"data_type"`
	ByteOrder  string    `json:"byte_order" csv:"byte_order" form:"byte_order"`
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
	Id         string    `json:"id" csv:"id" form:"id"`
	Name       string    `json:"name" csv:"name" form:"name"`
	Oid        string    `json:"oid" csv:"oid" form:"oid"`
	Unit       string    `json:"unit" csv:"unit" form:"unit"`
	ScriptId   string    `json:"script_id" csv:"script_id" form:"script_id"`
	DataFactor string    `json:"data_factor" csv:"data_factor" form:"data_factor"`
	Remark     string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" `
}

type ModbusVarTpl struct {
	Id            string `json:"id" csv:"id" form:"id"`
	Name          string `json:"name" csv:"name" form:"name"`
	DataType      string `json:"data_type" csv:"data_type" form:"data_type"`
	Unit          string `json:"unit" csv:"unit" form:"unit"`
	ByteOrder     string `json:"byte_order" csv:"byte_order" form:"byte_order"`
	ScriptId      string `json:"script_id" csv:"script_id" form:"script_id"`
	DataFactor    string `json:"data_factor" csv:"data_factor" form:"data_factor"`
	ChannelStatus string `json:"channel_status" csv:"channel_status" form:"channel_status"`
	Remark        string `json:"remark" csv:"remark" form:"remark"`
}

// ModbusCommand modbus 指令
type ModbusCommand struct {
	Id          string    `json:"id" csv:"id" form:"id"`
	Group       string    `json:"group" csv:"group" form:"group"` // 指令组
	Oid         string    `json:"oid" csv:"oid" form:"oid"`
	Name        string    `json:"name" csv:"name" form:"name"`                         // 指令名称
	Order       int       `json:"order" csv:"order" form:"order"`                      // 指令顺序
	DeviceId    string    `json:"device_id" csv:"device_id" form:"device_id"`          // 设备ID
	CommandType string    `json:"command_type" csv:"command_type" form:"command_type"` // 指令类型 data or script
	RegType     string    `json:"reg_type" csv:"reg_type" form:"reg_type"`             // 寄存器类型
	StartAddr   int       `json:"start_addr" csv:"start_addr" form:"start_addr"`       // 寄存器起始地址
	CommandData string    `json:"command_data" csv:"command_data" form:"command_data"` // 16进制字符串或者lua函数脚本
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	Oid      string `json:"oid" csv:"oid" form:"oid"`
	Register int    `json:"register" csv:"register" form:"register"`
	RegType  string `json:"reg_type" csv:"reg_type" form:"reg_type"`
	Length   int    `json:"length" csv:"length" form:"length"`
	Value    string `json:"value" csv:"value" form:"value"`
	Remark   string `json:"remark" csv:"remark" form:"remark"`
}

type MqttChannel struct {
	ID              string    `json:"id" csv:"id" form:"id"`
	Name            string    `json:"name" csv:"name" form:"name"`
	Server          string    `json:"server" csv:"server" form:"server"`
	ClientId        string    `json:"client_id" csv:"client_id" form:"client_id"`
	Username        string    `json:"username" csv:"username" form:"username"`
	Password        string    `json:"password" csv:"password" form:"password"`
	SubTopic        string    `json:"sub_topic" csv:"sub_topic" form:"sub_topic"`
	PubTopic        string    `json:"pub_topic" csv:"pub_topic" form:"pub_topic"`
	Will            string    `json:"will" csv:"will" form:"will"`
	Qos             int       `json:"qos" csv:"qos" form:"qos"`
	KeepAlive       int       `json:"keep_alive" csv:"keep_alive" form:"keep_alive"`
	PingTimeout     int       `json:"ping_timeout" csv:"ping_timeout" form:"ping_timeout"`
	RetryInterval   int       `json:"retry_interval" csv:"retry_interval" form:"retry_interval"`
	ClearSession    int       `json:"clear_session" csv:"clear_session" form:"clear_session"`
	ProtocolVersion int       `json:"protocol_version" csv:"protocol_version" form:"protocol_version"`
	Retained        int       `json:"retained" csv:"retained" form:"retained"`
	Debug           int       `json:"debug" csv:"debug" form:"debug"`
	Status          int       `json:"status" csv:"status" form:"status"`
	Remark          string    `json:"remark" csv:"remark" form:"remark"`
	LastBoot        time.Time `json:"last_boot" csv:"last_boot"`
	CreatedAt       time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TcpChannel struct {
	ID          string    `json:"id" csv:"id" form:"id"`
	Name        string    `json:"name" csv:"name" form:"name"`
	Server      string    `json:"server" csv:"server" form:"server"`
	ChannelType string    `json:"channel_type" csv:"channel_type" form:"channel_type"`
	Port        int       `json:"port" csv:"port" form:"port"`
	Timeout     int       `json:"timeout" csv:"timeout" form:"timeout"`
	Debug       int       `json:"debug" csv:"debug" form:"debug"`
	Status      int       `json:"status" csv:"status" form:"status"`
	Remark      string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt   time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HttpChannel struct {
	ID        string    `json:"id" csv:"id" form:"id"`
	Name      string    `json:"name" csv:"name" form:"name"`
	Url       string    `json:"url" csv:"url" form:"url"`
	Format    string    `json:"format" csv:"format" form:"format"`
	Header    string    `json:"header" csv:"header" form:"header"`
	Timeout   int       `json:"timeout" csv:"timeout" form:"timeout"`
	Debug     int       `json:"debug" csv:"debug" form:"debug"`
	Status    int       `json:"status" csv:"status" form:"status"`
	Remark    string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DataStream struct {
	ID          string    `json:"id" csv:"id" form:"id"`
	Name        string    `json:"name" csv:"name" form:"name"`
	SchedPolicy string    `json:"sched_policy" csv:"sched_policy" form:"sched_policy"` // 运行间隔
	MN          string    `json:"mn" csv:"mn" form:"mn"`                               // 关联设备 mn
	ScriptId    string    `json:"script_id" csv:"script_id" form:"script_id"`          // 数据脚本 ID
	MqttChids   string    `json:"mqtt_chids" csv:"mqtt_chids" form:"mqtt_chids"`       // MQTT 通道 ID
	TcpChids    string    `json:"tcp_chids" csv:"tcp_chids" form:"tcp_chids"`          // TCP 通道 ID
	HttpChids   string    `json:"http_chids" csv:"http_chids" form:"http_chids"`       // HTTP 通道 ID
	Remark      string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ControlStream struct {
	ID        string    `json:"id" csv:"id" form:"id"`
	Name      string    `json:"name" csv:"name" form:"name"`
	Event     string    `json:"event" csv:"event" form:"event"`
	CommandId string    `json:"command_id" csv:"command_id" form:"command_id"`
	Param     string    `json:"param" csv:"param" form:"param"`
	Remark    string    `json:"remark" csv:"remark" form:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Tables = []interface{}{
	&SysConfig{},
	&SysOpr{},
	&DataScript{},
	// hj212
	&Hj212Queue{},
	// modbus
	&ModbusDevice{},
	&ModbusReg{},
	&ModbusVar{},
	&ModbusVarTpl{},
	&ModbusCommand{},
	// iot device
	&IotDevice{},
	&DeviceRtdData{},
	&ModbusSlaveReg{},
	// channel
	&MqttChannel{},
	&TcpChannel{},
	&HttpChannel{},
	&DataStream{},
	&ControlStream{},
}
