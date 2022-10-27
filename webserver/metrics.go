package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/models"
)

type Metrics struct {
	Icon  string
	Value interface{}
	Title string
}

func NewMetrics(icon string, value interface{}, title string) *Metrics {
	return &Metrics{Icon: icon, Value: value, Title: title}
}

func (s *WebServer) initMetricsRouters() {
	s.get("/admin/metrics/system/cpuusage", s.MetricsCpuusage)
	s.get("/admin/metrics/system/memusage", s.MetricsMemusage)
	s.get("/admin/metrics/system/uptime", s.MetricsUptime)
	s.get("/admin/metrics/modbus/:name/count", s.ModbusCounter)
	s.get("/admin/metrics/modbus/line", s.ModbusMetricsLine)
	s.get("/admin/metrics/modbus/linedata", s.ModbusMetricsLineData)
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

func (s *WebServer) ModbusCounter(c echo.Context) error {
	var count int64
	var title string
	name := c.Param("name")
	switch name {
	case "device":
		title = "设备数"
		app.DB().Model(&models.ModbusDevice{}).Count(&count)
		return c.Render(http.StatusOK, "metrics",
			NewMetrics("mdi mdi-switch", count, title))
	case "slavereg":
		title = "从机寄存器"
		app.DB().Model(&models.ModbusSlaveReg{}).Count(&count)
	case "reg":
		title = "寄存器"
		app.DB().Model(&models.ModbusReg{}).Count(&count)
	}
	return c.Render(http.StatusOK, "metrics",
		NewMetrics("mdi mdi-link", count, title))
}

func (s *WebServer) ModbusMetricsLine(c echo.Context) error {
	return c.Render(http.StatusOK, "modbus_metrics_line", map[string]string{
		"mn": c.QueryParam("mn"),
	})
}

func (s *WebServer) ModbusMetricsLineData(c echo.Context) error {
	mn := c.QueryParam("mn")
	var dev models.ModbusDevice
	query := app.DB()
	if mn != "" {
		query = query.Where("mn = ?", mn)
	}
	if err := app.DB().First(&dev).Error; err != nil {
		common.Must(err)
	}

	type metricLineItem struct {
		Name   string          `json:"name"`
		Values [][]interface{} `json:"values"`
	}

	var regs []models.ModbusReg
	err := app.DB().Where("device_id = ?", dev.Id).Find(&regs).Error
	if err != nil {
		common.Must(err)
	}

	var items []metricLineItem
	for _, reg := range regs {
		item := metricLineItem{
			Name:   string(reg.Name),
			Values: make([][]interface{}, 0),
		}
		points, _ := app.TsDB().Select(
			fmt.Sprintf("modbus_metrics_%s_%s", dev.MN, reg.Id), nil,
			time.Now().Add(-8*time.Hour).Unix(), time.Now().Unix())
		for _, p := range points {
			item.Values = append(item.Values, []interface{}{p.Timestamp * 1000, p.Value})
		}
		items = append(items, item)
	}

	return c.JSON(200, map[string]interface{}{
		"title": "设备数据统计",
		"datas": items,
	})
}
