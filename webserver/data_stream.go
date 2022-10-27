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

func (s *WebServer) initDataStreamRouters() {
	s.get("/admin/datastream", s.DataStream)
	s.get("/admin/datastream/options", s.DataStreamOptions)
	s.get("/admin/datastream/query", s.DataStreamQuery)
	s.post("/admin/datastream/add", s.DataStreamAdd)
	s.post("/admin/datastream/update", s.DataStreamUpdate)
	s.post("/admin/datastream/save", s.DataStreamSave)
	s.get("/admin/datastream/delete", s.DataStreamDelete)

}

func (s *WebServer) DataStream(c echo.Context) error {
	return c.Render(http.StatusOK, "data_stream", map[string]string{})
}

func (s *WebServer) DataStreamOptions(c echo.Context) error {
	var data []models.DataStream
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

func (s *WebServer) DataStreamQuery(c echo.Context) error {
	var data []models.DataStream
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) DataStreamAdd(c echo.Context) error {
	form := new(models.DataStream)
	common.Must(c.Bind(form))
	form.ID = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) DataStreamUpdate(c echo.Context) error {
	form := new(models.DataStream)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) DataStreamSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.DataStream)
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
		common.Must(app.DB().Delete(models.DataStream{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) DataStreamDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.DataStream{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
