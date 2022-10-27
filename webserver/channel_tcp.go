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

func (s *WebServer) initTcpChannelRouters() {
	s.get("/admin/channel/tcp", s.TcpChannel)
	s.get("/admin/channel/tcp/options", s.TcpChannelOptions)
	s.get("/admin/channel/tcp/query", s.TcpChannelQuery)
	s.post("/admin/channel/tcp/add", s.TcpChannelAdd)
	s.post("/admin/channel/tcp/update", s.TcpChannelUpdate)
	s.post("/admin/channel/tcp/save", s.TcpChannelSave)
	s.get("/admin/channel/tcp/delete", s.TcpChannelDelete)

}

func (s *WebServer) TcpChannel(c echo.Context) error {
	return c.Render(http.StatusOK, "tcp_channel", map[string]string{})
}

func (s *WebServer) TcpChannelOptions(c echo.Context) error {
	var data []models.TcpChannel
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

func (s *WebServer) TcpChannelQuery(c echo.Context) error {
	var data []models.TcpChannel
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) TcpChannelAdd(c echo.Context) error {
	form := new(models.TcpChannel)
	common.Must(c.Bind(form))
	form.ID = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("url", form.Server)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) TcpChannelUpdate(c echo.Context) error {
	form := new(models.TcpChannel)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("url", form.Server)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) TcpChannelSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.TcpChannel)
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
		common.Must(app.DB().Delete(models.TcpChannel{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) TcpChannelDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.TcpChannel{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
