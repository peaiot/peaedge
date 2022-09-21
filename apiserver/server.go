package apiserver

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/log"
)

var server *ApiServer

type ApiServer struct {
	root *echo.Echo
}

func Listen() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	server = NewApiServer()
	server.initRouter()
	return server.Start()
}

func NewApiServer() *ApiServer {
	s := new(ApiServer)
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
	p := prometheus.NewPrometheus("peaedge", nil)
	p.Use(s.root)
	s.root.HideBanner = true
	s.root.Logger.SetLevel(common.If(app.Config().Web.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	s.root.Debug = app.Config().Web.Debug
	return s
}

// Start 启动服务器
func (s *ApiServer) Start() error {
	log.Infof("启动 API 服务器 %s:%d", app.Config().Web.Host, app.Config().Web.Port)
	err := s.root.Start(fmt.Sprintf("%s:%d", app.Config().Web.Host, app.Config().Web.Port))
	if err != nil {
		log.Errorf("启动 API 服务器错误 %s", err.Error())
	}
	return err
}
