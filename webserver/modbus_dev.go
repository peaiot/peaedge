package webserver

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/web"
	"github.com/toughstruct/peaedge/models"
)

func (s *WebServer) initModbusDevRouters() {
	// modbus
	s.get("/admin/modbus/proto/options", s.ModbusProtoOptions)
	s.get("/admin/modbus/device", s.ModbusDevice)
	s.get("/admin/modbus/device/query", s.ModbusDeviceQuery)
	s.post("/admin/modbus/device/save", s.ModbusDeviceSave)
	s.get("/admin/modbus/device/delete", s.ModbusDeviceDelete)
}

func (s *WebServer) ModbusDevice(c echo.Context) error {
	return c.Render(http.StatusOK, "modbus_device", map[string]string{})
}

func (s *WebServer) ModbusDeviceQuery(c echo.Context) error {
	var data []models.ModbusDevice
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ModbusDeviceSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ModbusDevice)
	common.Must(c.Bind(form))
	switch op {
	case "insert":
		form.Id = common.UUIDBase32()
		common.Must(app.DB().Create(form).Error)
		return c.JSON(200, map[string]interface{}{"id": form.Id})
	case "update":
		common.Must(app.DB().Updates(form).Error)
		return c.JSON(200, map[string]interface{}{"status": "updated"})
	case "delete":
		common.Must(app.DB().Delete(models.ModbusDevice{}, form.Id).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ModbusDeviceDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ModbusDevice{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
