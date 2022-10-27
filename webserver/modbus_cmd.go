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

func (s *WebServer) initModbusCommandRouters() {
	s.get("/admin/modbus/command", s.ModbusCommand)
	s.get("/admin/modbus/command/options", s.ModbusCommandOptions)
	s.get("/admin/modbus/command/query", s.ModbusCommandQuery)
	s.post("/admin/modbus/command/add", s.ModbusCommandAdd)
	s.post("/admin/modbus/command/update", s.ModbusCommandUpdate)
	s.post("/admin/modbus/command/save", s.ModbusCommandSave)
	s.get("/admin/modbus/command/delete", s.ModbusCommandDelete)

}

func (s *WebServer) ModbusCommand(c echo.Context) error {
	return c.Render(http.StatusOK, "modbus_command", map[string]string{})
}

func (s *WebServer) ModbusCommandOptions(c echo.Context) error {
	var data []models.ModbusCommand
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

func (s *WebServer) ModbusCommandQuery(c echo.Context) error {
	var data []models.ModbusCommand
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ModbusCommandAdd(c echo.Context) error {
	form := new(models.ModbusCommand)
	common.Must(c.Bind(form))
	form.Id = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) ModbusCommandUpdate(c echo.Context) error {
	form := new(models.ModbusCommand)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) ModbusCommandSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ModbusCommand)
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
		common.Must(app.DB().Delete(models.ModbusCommand{}, form.Id).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ModbusCommandDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ModbusCommand{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
