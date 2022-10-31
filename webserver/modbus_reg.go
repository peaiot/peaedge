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

func (s *WebServer) initModbusRegRouters() {
	// modbus reg
	s.get("/admin/modbus/reg/query", s.ModbusRegQuery)
	s.post("/admin/modbus/reg/save", s.ModbusRegSave)
	s.get("/admin/modbus/reg/delete", s.ModbusRegDelete)
}

func (s *WebServer) ModbusRegQuery(c echo.Context) error {
	devId := c.QueryParam("device_id")
	var data []models.ModbusReg
	query := app.DB()
	if devId != "" {
		query = query.Where("device_id = ?", devId)
	}
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ModbusRegSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ModbusReg)
	common.Must(c.Bind(form))
	switch op {
	case "insert":
		form.Id = common.UUIDBase32()
		common.Must(app.DB().Create(&form).Error)
		return c.JSON(200, map[string]interface{}{"id": form.Id})
	case "update":
		common.Must(app.DB().Select(
			"name", "reg_type", "start_addr", "oid", "byte_order", "data_type",
			"data_len", "access_type", "var_id", "status",
		).Updates(&form).Error)
		return c.JSON(200, map[string]interface{}{"status": "updated"})
	case "delete":
		common.Must(app.DB().Delete(models.ModbusReg{}, form.Id).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ModbusRegDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ModbusReg{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
