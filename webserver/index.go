package webserver

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/models"
)

func (s *WebServer) Index(c echo.Context) error {
	sess, _ := session.Get(UserSession, c)
	username := sess.Values[UserSessionName]
	if username == nil || username == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login?errmsg=用户未登录或登录过期")
	}
	return c.Render(http.StatusOK, "index", map[string]interface{}{})
}

func (s *WebServer) Dashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", map[string]string{})
}

func (s *WebServer) Menus(c echo.Context) error {
	sess, _ := session.Get(UserSession, c)
	switch sess.Values[UserSessionLevel] {
	case "super":
		return c.JSONBlob(http.StatusOK, assets.AdminMenudata)
	case "opr":
		return c.JSONBlob(http.StatusOK, assets.OprMenudata)
	default:
		return c.JSONBlob(http.StatusOK, nil)
	}
}

func (s *WebServer) Login(c echo.Context) error {
	errmsg := c.QueryParam("errmsg")
	return c.Render(http.StatusOK, "login", map[string]interface{}{
		"errmsg":    errmsg,
		"LoginLogo": "/static/images/login-logo.png",
	})
}

func (s *WebServer) Logout(c echo.Context) error {
	sess, _ := session.Get(UserSession, c)
	sess.Values = make(map[interface{}]interface{})
	_ = sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

func (s *WebServer) LoginPost(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" || password == "" {
		return c.Redirect(http.StatusMovedPermanently, "/login?errmsg=用户名和密码不能为空")
	}
	var user models.SysOpr
	err := app.DB().Where("username=?", username).First(&user).Error
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, "/login?errmsg=用户不存在")
	}

	if common.Sha256Hash(password) != user.Password {
		return c.Redirect(http.StatusMovedPermanently, "/login?errmsg=密码错误")
	}

	sess, _ := session.Get(UserSession, c)
	sess.Values[UserSessionName] = username
	sess.Values[UserSessionLevel] = user.Level
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return echo.NewHTTPError(http.StatusMovedPermanently, err.Error())
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}
