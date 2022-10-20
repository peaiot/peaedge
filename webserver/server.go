package webserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/tpl"
	"github.com/toughstruct/peaedge/log"
)

var server *WebServer

type WebServer struct {
	root *echo.Echo
}

const UserSession = "user_session"
const UserSessionName = "user_session_name"
const UserSessionLevel = "user_session_level"
const ConstCookieName = "teamsacs_cookie"

func Listen() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	server = NewWebServer()
	server.initRouters()
	return server.Start()
}

func NewWebServer() *WebServer {
	s := new(WebServer)
	s.root = echo.New()
	s.root.Pre(middleware.RemoveTrailingSlash())
	s.root.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
				}
			}()
			return next(c)
		}
	})
	s.root.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "freeradius ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${latency} ${bytes_in} ${bytes_out} ${error}\n",
		Output: os.Stdout,
	}))
	s.root.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(assets.StaticFs))))
	s.root.Renderer = tpl.NewCommonTemplate(assets.TemplatesFs, []string{"templates"}, map[string]interface{}{
		"pagever": func() int64 {
			return time.Now().Unix()
		},
		"buildver": func() string {
			return strings.TrimSpace(assets.BuildVer)
		},
	})
	p := prometheus.NewPrometheus(app.Config().System.Appid, nil)
	p.Use(s.root)

	// session 中间件， 采用 Cookie 存储方式
	sessStore := sessions.NewCookieStore([]byte(app.Config().Web.Secret))
	sessStore.MaxAge(3600 * 2)
	s.root.Use(session.Middleware(sessStore))
	s.root.Use(sessionCheck())

	s.root.HideBanner = true
	s.root.Logger.SetLevel(common.If(app.Config().Web.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	s.root.Debug = app.Config().Web.Debug
	return s
}

// Start 启动服务器
func (s *WebServer) Start() error {
	log.Infof("启动 WEB 服务器 %s:%d", app.Config().Web.Host, app.Config().Web.Port)
	err := server.root.Start(fmt.Sprintf("%s:%d", app.Config().Web.Host, app.Config().Web.Port))
	if err != nil {
		log.Errorf("启动 WEB 服务器错误 %s", err.Error())
	}
	return err
}

// 检查 Session
func sessionCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			var skips = []string{
				"/login",
				"/logout",
				"/admin/login",
				"/static",
			}

			for _, prefix := range skips {
				if strings.HasPrefix(c.Path(), prefix) {
					return next(c)
				}
			}

			sess, _ := session.Get(UserSession, c)
			username := sess.Values[UserSessionName]
			if username == nil || username == "" {
				return c.Redirect(http.StatusTemporaryRedirect, "/login?errmsg=用户未登录或登录过期")
			}
			return next(c)
		}
	}
}
