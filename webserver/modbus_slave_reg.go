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

func (s *WebServer) initModbusSlaveRegRouters() {
	s.get("/admin/modbus/slavereg", s.ModbusSlaveReg)
	s.get("/admin/modbus/slavereg/options", s.ModbusSlaveRegOptions)
	s.get("/admin/modbus/slavereg/query", s.ModbusSlaveRegQuery)
	s.post("/admin/modbus/slavereg/save", s.ModbusSlaveRegSave)
	s.get("/admin/modbus/slavereg/delete", s.ModbusSlaveRegDelete)

}

func (s *WebServer) ModbusSlaveReg(c echo.Context) error {
	return c.Render(http.StatusOK, "modbus_slave_reg", map[string]string{})
}

func (s *WebServer) ModbusSlaveRegOptions(c echo.Context) error {
	var data []models.ModbusSlaveReg
	err := app.DB().Find(&data).Error
	common.Must(err)
	var options []*web.JsonOptions
	options = append(options, &web.JsonOptions{
		Id:    common.NA,
		Value: "",
	})
	for _, d := range data {
		options = append(options, &web.JsonOptions{
			Id:    d.ID,
			Value: d.Name,
		})
	}
	return c.JSON(http.StatusOK, options)
}

func (s *WebServer) ModbusSlaveRegQuery(c echo.Context) error {
	var data []models.ModbusSlaveReg
	query := app.DB()
	err := query.Order("register asc").Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ModbusSlaveRegSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ModbusSlaveReg)
	common.Must(c.Bind(form))
	switch op {
	case "insert":
		form.ID = common.UUID()
		common.Must(app.DB().Create(form).Error)
		return c.JSON(200, map[string]interface{}{"id": form.ID})
	case "update":
		common.Must(app.DB().Updates(form).Error)
		return c.JSON(200, map[string]interface{}{"status": "updated"})
	case "delete":
		common.Must(app.DB().Delete(models.ModbusSlaveReg{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ModbusSlaveRegDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ModbusSlaveReg{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
