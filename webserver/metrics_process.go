package webserver

import (
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/app"
)

type processMetricLineItem struct {
	Id    int     `json:"id"`
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

func (s *WebServer) MonitorCpuMetricsLineData(c echo.Context) error {

	var items []processMetricLineItem

	points, _ := app.TsDB().Select("peaedge_cpuuse", nil,
		time.Now().Add(-1*time.Hour).Unix(), time.Now().Unix())
	for i, p := range points {
		items = append(items, processMetricLineItem{
			Id:    i + 1,
			Time:  time.Unix(p.Timestamp, 0).Format("15:04"),
			Value: p.Value,
		})
	}

	result := statProcessMetricLine(items)

	return c.JSON(200, result)
}

func (s *WebServer) MonitorMemMetricsLineData(c echo.Context) error {

	var items []processMetricLineItem

	points, _ := app.TsDB().Select("peaedge_memuse", nil,
		time.Now().Add(-1*time.Hour).Unix(), time.Now().Unix())
	for i, p := range points {
		items = append(items, processMetricLineItem{
			Id:    i + 1,
			Time:  time.Unix(p.Timestamp, 0).Format("15:04"),
			Value: p.Value,
		})
	}

	result := statProcessMetricLine(items)
	return c.JSON(200, result)
}

func statProcessMetricLine(src []processMetricLineItem) []processMetricLineItem {
	df := dataframe.LoadStructs(src)
	groups := df.GroupBy("Time")
	aggre := groups.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_MEAN}, []string{"Value"})
	sorted := aggre.Arrange(
		dataframe.Sort("Time"), // Sort in ascending order
	)
	var nitems []processMetricLineItem
	for i, vals := range sorted.Records() {
		if i == 0 {
			continue
		}
		nitems = append(nitems, processMetricLineItem{
			Id:    i,
			Time:  vals[0],
			Value: cast.ToFloat64(vals[1]),
		})
	}
	return nitems
}
