package tests

import (
	"testing"

	"github.com/go-gota/gota/dataframe"
	"github.com/spf13/cast"
)

func TestAvgGroup(t *testing.T) {
	type metricLineItem struct {
		Id    int     `json:"id"`
		Time  string  `json:"time"`
		Value float64 `json:"value"`
	}

	var items = []metricLineItem{
		{1, "15:04", 1.0},
		{2, "15:04", 2.0},
		{3, "15:04", 3.0},
		{4, "15:05", 4.0},
		{5, "15:05", 5.0},
		{6, "15:06", 6.0},
		{7, "15:06", 7.0},
		{8, "15:07", 8.0},
		{9, "15:07", 9.0},
	}

	df := dataframe.LoadStructs(items)
	t.Log(df)
	groups := df.GroupBy("Time")
	t.Log(groups)
	aggre := groups.Aggregation([]dataframe.AggregationType{dataframe.Aggregation_MEAN}, []string{"Value"})
	t.Log(aggre)
	var nitems []metricLineItem
	for i, vals := range aggre.Records() {
		if i == 0 {
			continue
		}
		nitems = append(nitems, metricLineItem{
			Id:    i,
			Time:  vals[0],
			Value: cast.ToFloat64(vals[1]),
		})
	}
	t.Log(nitems)

}
