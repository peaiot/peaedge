package webserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/web"
	"github.com/toughstruct/peaedge/golua"
	"github.com/toughstruct/peaedge/models"
)

func (s *WebServer) initDataScriptRouters() {
	s.get("/admin/datascript", s.DataScript)
	s.get("/admin/datascript/options", s.DataScriptOptions)
	s.get("/admin/datascript/:funcname/options", s.DataScriptOptions)
	s.get("/admin/datascript/query", s.DataScriptQuery)
	s.post("/admin/datascript/save", s.DataScriptSave)
	s.post("/admin/datascript/add", s.DataScriptAdd)
	s.post("/admin/datascript/update", s.DataScriptUpdate)
	s.post("/admin/datascript/test", s.DataScriptTest)
	s.get("/admin/datascript/delete", s.DataScriptDelete)

}

func (s *WebServer) DataScript(c echo.Context) error {
	return c.Render(http.StatusOK, "data_script", map[string]string{})
}

func (s *WebServer) DataScriptOptions(c echo.Context) error {
	stype := c.Param("funcname")
	var data []models.DataScript
	query := app.DB()
	if stype != "" {
		query = query.Where("func_name = ?", stype)
	}
	err := query.Find(&data).Error
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

func (s *WebServer) DataScriptQuery(c echo.Context) error {
	var data []models.DataScript
	query := app.DB()
	err := query.Find(&data).Error
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (s *WebServer) DataScriptAdd(c echo.Context) error {
	form := new(models.DataScript)
	common.Must(c.Bind(form))
	form.ID = common.UUIDBase32()
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("func_name", form.FuncName)
	common.Must(app.DB().Create(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) DataScriptUpdate(c echo.Context) error {
	form := new(models.DataScript)
	common.Must(c.Bind(form))
	common.CheckEmpty("name", form.Name)
	common.CheckEmpty("func_name", form.FuncName)
	common.Must(app.DB().Save(form).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

func (s *WebServer) DataScriptTest(c echo.Context) error {
	args := c.FormValue("args")
	src := c.FormValue("src")
	fname := c.FormValue("func")
	switch fname {
	case "HandlerModbusRtd":
		ret, err := golua.HandlerModbusRtd(src, cast.ToInt(args))
		if err != nil {
			return c.JSON(http.StatusOK, web.RestError(err.Error()))
		}
		return c.JSON(http.StatusOK, web.RestSucc("Result = "+cast.ToString(ret)))
	case "HandlerDataStream":
		ret, err := golua.HandlerDataStream(src, args)
		if err != nil {
			return c.JSON(http.StatusOK, web.RestError(err.Error()))
		}
		return c.JSON(http.StatusOK, web.RestSucc("Result = "+cast.ToString(ret)))
	default:
		return c.JSON(http.StatusOK, web.RestError(fmt.Sprintf("func %s not support", fname)))
	}
}

func (s *WebServer) DataScriptSave(c echo.Context) error {
	op := c.FormValue("webix_operation")
	form := new(models.DataScript)
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
		common.Must(app.DB().Delete(models.DataScript{}, form.ID).Error)
		return c.JSON(200, make(map[string]interface{}))
	default:
		return c.JSON(http.StatusOK, make(map[string]interface{}))
	}

}

func (s *WebServer) DataScriptDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	common.Must(app.DB().Delete(models.DataScript{}, strings.Split(ids, ",")).Error)
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
