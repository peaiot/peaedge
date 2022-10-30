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

func (s *WebServer) initControlStreamRouters() {
	s.get("/admin/controlstream", s.ControlStream)
	s.get("/admin/controlstream/options", s.ControlStreamOptions)
	s.get("/admin/controlstream/query", s.ControlStreamQuery)
	s.post("/admin/controlstream/add", s.ControlStreamAdd)
	s.post("/admin/controlstream/update", s.ControlStreamUpdate)
	s.post("/admin/controlstream/save", s.ControlStreamSave)
	s.get("/admin/controlstream/delete", s.ControlStreamDelete)
}

func (s *WebServer) ControlStream(c echo.Context) error {
	return c.Render(http.StatusOK, "control_stream", map[string]string{})
}

func (s *WebServer) ControlStreamOptions(c echo.Context) error {
	var data []models.ControlStream
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

func (s *WebServer) ControlStreamQuery(c echo.Context) error {
	var data []models.ControlStream
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) ControlStreamAdd(c echo.Context) error {
	form := new(models.ControlStream)
	common.Must(c.Bind(form))
	form.ID = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) ControlStreamUpdate(c echo.Context) error {
	form := new(models.ControlStream)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) ControlStreamSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.ControlStream)
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
		common.Must(app.DB().Delete(models.ControlStream{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) ControlStreamDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.ControlStream{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
