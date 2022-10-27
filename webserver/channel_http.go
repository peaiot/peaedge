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

func (s *WebServer) initHttpChannelRouters() {
	s.get("/admin/channel/http", s.HttpChannel)
	s.get("/admin/channel/http/options", s.HttpChannelOptions)
	s.get("/admin/channel/http/query", s.HttpChannelQuery)
	s.post("/admin/channel/http/add", s.HttpChannelAdd)
	s.post("/admin/channel/http/update", s.HttpChannelUpdate)
	s.post("/admin/channel/http/save", s.HttpChannelSave)
	s.get("/admin/channel/http/delete", s.HttpChannelDelete)

}

func (s *WebServer) HttpChannel(c echo.Context) error {
	return c.Render(http.StatusOK, "http_channel", map[string]string{})
}

func (s *WebServer) HttpChannelOptions(c echo.Context) error {
	var data []models.HttpChannel
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

func (s *WebServer) HttpChannelQuery(c echo.Context) error {
	var data []models.HttpChannel
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) HttpChannelAdd(c echo.Context) error {
	form := new(models.HttpChannel)
	form.ID = common.UUIDBase32()
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("url", form.Url)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) HttpChannelUpdate(c echo.Context) error {
	form := new(models.HttpChannel)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("url", form.Url)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) HttpChannelSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.HttpChannel)
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
		common.Must(app.DB().Delete(models.HttpChannel{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) HttpChannelDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.HttpChannel{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
