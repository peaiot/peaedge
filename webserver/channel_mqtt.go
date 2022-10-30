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

func (s *WebServer) initMqttChannelRouters() {
	s.get("/admin/channel/mqtt", s.MqttChannel)
	s.get("/admin/channel/mqtt/options", s.MqttChannelOptions)
	s.get("/admin/channel/mqtt/query", s.MqttChannelQuery)
	s.post("/admin/channel/mqtt/add", s.MqttChannelAdd)
	s.post("/admin/channel/mqtt/update", s.MqttChannelUpdate)
	s.post("/admin/channel/mqtt/save", s.MqttChannelSave)
	s.get("/admin/channel/mqtt/delete", s.MqttChannelDelete)

}

func (s *WebServer) MqttChannel(c echo.Context) error {
	return c.Render(http.StatusOK, "mqtt_channel", map[string]string{})
}

func (s *WebServer) MqttChannelOptions(c echo.Context) error {
	var data []models.MqttChannel
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

func (s *WebServer) MqttChannelQuery(c echo.Context) error {
	var data []models.MqttChannel
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) MqttChannelAdd(c echo.Context) error {
	form := new(models.MqttChannel)
	common.Must(c.Bind(form))
	form.ID = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("server", form.Server)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) MqttChannelUpdate(c echo.Context) error {
	form := new(models.MqttChannel)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("server", form.Server)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) MqttChannelSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.MqttChannel)
	common.Must(c.Bind(form))
	switch op {
	case "insert":
		form.ID = common.UUIDBase32()
		common.Must(app.DB().Create(form).Error)
		return c.JSON(200, map[string]interface{}{"id": form.ID})
	case "update":
		common.Must(app.DB().Updates(form).Error)
		return c.JSON(200, map[string]interface{}{"status": "updated"})
	case "delete":
		common.Must(app.DB().Delete(models.MqttChannel{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) MqttChannelDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.MqttChannel{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
