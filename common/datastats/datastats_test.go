package datastats

import (
	_ "embed"
	"os"
	"testing"

	"github.com/go-gota/gota/dataframe"
)

//go:embed testdata.csv
var testcsv string

func TestStatdf(t *testing.T) {
	file, err := os.Open("testdata.csv")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	df := dataframe.ReadCSV(file)
	t.Log(df)

	// ff := df.Filter(
	// 	dataframe.F{
	// 		Colname:    "value",
	// 		Comparator: series.Neq,
	// 		Comparando: "N/A",
	// 	},
	// )
	// if ff.Nrow() == 0 {
	// 	t.Fatal("filter failed")
	// }
	// ff := df.Rapply(func(s series.Series) series.Series {
	// 	s.Elem(0).Set(s.Val(0).(string)[:13])
	// 	return s
	// })
	// t.Log(ff)

	df = df.Arrange(
		dataframe.Sort("datetime"), // Sort in ascending order
		dataframe.Sort("name"),     // Sort in ascending order
	)

	t.Log(df.Describe())

	egg := df.GroupBy("datetime", "name").
		Aggregation([]dataframe.AggregationType{dataframe.Aggregation_SUM,
			dataframe.Aggregation_MEAN,
			dataframe.Aggregation_MAX}, []string{"value", "value", "value"})

	t.Log(egg)

}
