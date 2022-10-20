package webserver

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/web"
	"github.com/toughstruct/peaedge/models"
)

func (s *WebServer) initOprRouters() {
	// modbus
	s.get("/admin/opr/options", s.SysOprOptions)
	s.get("/admin/opr", s.SysOpr)
	s.get("/admin/opr/query", s.SysOprQuery)
	s.post("/admin/opr/save", s.SysOprSave)
	s.post("/admin/opr/add", s.SysOprAdd)
	s.post("/admin/opr/update", s.SysOprUpdate)
	s.get("/admin/opr/delete", s.SysOprDelete)
}

func (s *WebServer) SysOpr(c echo.Context) error {
	return c.Render(http.StatusOK, "opr", map[string]string{})
}

func (s *WebServer) SysOprOptions(c echo.Context) error {
	var data []models.SysOpr
	err := app.DB().Find(&data).Error
	common.Must(err)
	var options []*web.JsonOptions
	options = append(options, &web.JsonOptions{
		Id:    common.NA,
		Value: "",
	})
	for _, d := range data {
		options = append(options, &web.JsonOptions{
			Id:    cast.ToString(d.ID),
			Value: d.Realname,
		})
	}
	return c.JSON(http.StatusOK, options)
}

func (s *WebServer) SysOprQuery(c echo.Context) error {
	var data []models.SysOpr
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) SysOprAdd(c echo.Context) error {
	form := new(models.SysOpr)
	form.ID = common.UUIDBase32()
	common.Must(c.Bind(form))
	common.CheckEmpty("Username", form.Username)
	common.CheckEmpty("Realname", form.Realname)
	common.CheckEmpty("Password", form.Password)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) SysOprUpdate(c echo.Context) error {
	form := new(models.SysOpr)
	common.Must(c.Bind(form))
	common.CheckEmpty("Username", form.Username)
	common.CheckEmpty("Realname", form.Realname)
	if form.Password != "" {
		form.Password = common.Sha256Hash(form.Password)
	}
	common.Must(app.DB().Updates(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) SysOprSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.SysOpr)
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
		common.Must(app.DB().Delete(models.SysOpr{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) SysOprDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.SysOpr{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
