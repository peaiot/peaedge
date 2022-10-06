package webserver

import (
	"github.com/labstack/echo/v4"
	"github.com/toughstruct/peaedge/common/log"
)

func (s *WebServer) initRouters() {
	s.get("/", s.Index)
	s.get("/login", s.Login)
	s.get("/logout", s.Logout)
	s.post("/login", s.LoginPost)
	s.get("/admin/menu.json", s.Menus)
	s.get("/admin/dashboard", s.Dashboard)

	// metrics
	s.get("/admin/metrics/system/cpuusage", s.MetricsCpuusage)
	s.get("/admin/metrics/system/memusage", s.MetricsMemusage)
	s.get("/admin/metrics/system/uptime", s.MetricsUptime)

	s.initOptionsRouters()
	s.initModbusDevRouters()
	s.initModbusRegRouters()
	s.initModbusVarRouters()
	s.initDataScriptRouters()
	s.initOprRouters()
}

func (s *WebServer) get(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	log.Debugf("Add GET Router %s", path)
	return s.root.GET(path, h, m...)
}

func (s *WebServer) post(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	log.Debugf("Add POST Router %s", path)
	return s.root.POST(path, h, m...)
}
