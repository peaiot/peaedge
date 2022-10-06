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
	s.get("/admin/regtype/options", s.RegisterTypeOptions)
	s.get("/admin/datatype/options", s.DatatypeOptions)
	s.get("/admin/device/ports", s.PortOptions)
	s.get("/admin/byteorder/options", s.ByteOrderOptions)
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
func (s *WebServer) RegisterTypeOptions(c echo.Context) error {
	return c.JSON(200, []web.JsonOptions{
		{Id: "InputRegister", Value: "InputRegister"},
		{Id: "HoldingRegister", Value: "HoldingRegister"},
		{Id: "DiscreteInput", Value: "DiscreteInput"},
		{Id: "Coil", Value: "Coil"},
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
