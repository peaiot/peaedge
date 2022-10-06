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

func (s *WebServer) initModbusVarRouters() {
	s.get("/admin/modbus/var", s.ModbusVar)
	s.get("/admin/modbus/var/options", s.ModbusVarOptions)
	s.get("/admin/modbus/var/query", s.ModbusVarQuery)
	s.post("/admin/modbus/var/save", s.ModbusVarSave)
	s.get("/admin/modbus/var/delete", s.ModbusVarDelete)

}

func (s *WebServer) ModbusVar(c echo.Context) error {
	return c.Render(http.StatusOK, "modbus_var", map[string]string{})
}

func (s *WebServer) ModbusVarOptions(c echo.Context) error {
	var data []models.ModbusVar
	err := app.DB().Find(&data).Error
	common.Must(err)
	var options []*web.JsonOptions
	options = append(options, &web.JsonOptions{
		Id:    common.NA,
		Value: "",
	})
	for _, d := range data {
		options = append(options, &web.JsonOptions{
			Id:    d.Id,
			Value: d.Name,
		})
	}
	return c.JSON(http.StatusOK, options)
}

func (s *WebServer) ModbusVarQuery(c echo.Context) error {
	var data []models.ModbusVar
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ModbusVarSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ModbusVar)
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
		common.Must(app.DB().Delete(models.ModbusVar{}, form.Id).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ModbusVarDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ModbusVar{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
