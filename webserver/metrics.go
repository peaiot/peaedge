package webserver

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Metrics struct {
	Icon  string
	Value interface{}
	Title string
}

func NewMetrics(icon string, value interface{}, title string) *Metrics {
	return &Metrics{Icon: icon, Value: value, Title: title}
}

// MetricsCpuusage /admin/metrics/system/cpuusage
func (s *WebServer) MetricsCpuusage(c echo.Context) error {
	_cpuuse, err := cpu.Percent(0, false)
	_cpucount, _ := cpu.Counts(false)
	if err != nil {
		_cpuuse = []float64{0}
	}
	return c.Render(http.StatusOK, "metrics",
		NewMetrics("mdi mdi-circle-slice-2",
			fmt.Sprintf("%.2f%%", _cpuuse[0]),
			fmt.Sprintf("Cpu %d Core", _cpucount)))
}

// MetricsMemusage /admin/metrics/system/memusage
func (s *WebServer) MetricsMemusage(c echo.Context) error {
	_meminfo, err := mem.VirtualMemory()
	_usage := 0.0
	_total := uint64(0)
	if err == nil {
		_usage = _meminfo.UsedPercent
		_total = _meminfo.Total / (1000 * 1000 * 1000)
	}
	return c.Render(http.StatusOK, "metrics",
		NewMetrics("mdi mdi-memory", fmt.Sprintf("%.2f%%", _usage),
			fmt.Sprintf("Memory Total: %d G", _total)))
}

// MetricsUptime /admin/metrics/system/uptime
func (s *WebServer) MetricsUptime(c echo.Context) error {
	hinfo, err := host.Info()
	_hour := uint64(0)
	if err == nil {
		_hour = hinfo.Uptime
	}
	return c.Render(http.StatusOK, "metrics",
		NewMetrics("mdi mdi-clock",
			fmt.Sprintf("%.1f Hour",
				float64(_hour)/float64(3600)), "运行时长"))
}
