package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/toughstruct/peaedge/common/modbus"
	"github.com/toughstruct/peaedge/common/web"
)

func (s *WebServer) initOptionsRouters() {
	// options
	s.get("/admin/status/options", s.StatusOptions)
	s.get("/admin/luafunc/options", s.LuaFuncOptions)
	s.get("/admin/regtype/options", s.RegisterTypeOptions)
	s.get("/admin/datatype/options", s.DatatypeOptions)
	s.get("/admin/device/ports", s.PortOptions)
	s.get("/admin/byteorder/options", s.ByteOrderOptions)
	s.get("/admin/sched/options", s.SchedOptions)
}

func (s *WebServer) ModbusProtoOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "modbusrtu", Value: "ModbusRTU"},
		{Id: "modbustcp", Value: "ModbusTCP"},
	})
}

func (s *WebServer) StatusOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "enabled", Value: "启用"},
		{Id: "disabled", Value: "停用"},
	})
}

func (s *WebServer) LuaFuncOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "HandlerModbusRtd", Value: "HandlerModbusRtd - 常规 Modbus 实时值计算"},
		{Id: "HandlerModbusCommand", Value: "HandlerModbusCommand - Modbus 操作指令数据处理"},
		{Id: "HandlerMixedRegister", Value: "HandlerMixedRegister - 混合寄存器计算"},
		{Id: "HandlerDataStream", Value: "HandlerDataStream - 数据流封装处理"},
	})
}

func (s *WebServer) RegisterTypeOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "InputRegister", Value: "InputRegister"},
		{Id: "HoldingRegister", Value: "HoldingRegister"},
		{Id: "DiscreteInput", Value: "DiscreteInput"},
		{Id: "Coil", Value: "Coil"},
		{Id: "MixedRegister", Value: "MixedRegister"},
	})
}

func (s *WebServer) DatatypeOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "int16", Value: "int16"},
		{Id: "uint16", Value: "uint16"},
		{Id: "float", Value: "float"},
		{Id: "float64", Value: "float64"},
	})
}
func (s *WebServer) ByteOrderOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "BigEndian", Value: "BigEndian-ABCD"},
		{Id: "BigEndianSwap", Value: "BigEndianSwap-BADC"},
		{Id: "LittleEndian", Value: "LittleEndian-DCBA"},
		{Id: "LittleEndianSwap", Value: "LittleEndianSwap-CDAB"},
	})
}

func (s *WebServer) PortOptions(c echo.Context) error {
	ports := modbus.ScanPorts()
	return c.JSON(http.StatusOK, ports)
}

func (s *WebServer) SchedOptions(c echo.Context) error {
	var opts = make([]web.JsonOptions, 0)
	opts = append(opts, web.JsonOptions{Id: "minute", Value: "每1分钟执行一次"})
	opts = append(opts, web.JsonOptions{Id: "5minute", Value: "每5分钟执行一次"})
	opts = append(opts, web.JsonOptions{Id: "10minute", Value: "每10分钟执行一次"})
	opts = append(opts, web.JsonOptions{Id: "hour", Value: "每小时执行一次"})
	opts = append(opts, web.JsonOptions{Id: "daily", Value: "每日执行一次"})
	return c.JSON(http.StatusOK, opts)
}
