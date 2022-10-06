package goscript

import (
	"testing"
)

func TestRunGopScript(t *testing.T) {
	s := `import "math"

// HandlerModbusData 
// 计算传感器值 计算公式 = （（采集值-最小取样值）/（最大取样值-最小取样值）*（最大量程-最小量程） ）  + 最小量程
func HandlerModbusData(src float64) float64 {
	var minSpVal float64 = 819.0 // 最小取样值 如 819
	var maxSpVal float64 = 4095.0 // 最大取样值 如 4095
	var minRangeVal float64 = 0.0 //最小量程 0
	var maxRangeVal float64 = 300 //最大量程 300
	if src < minSpVal {
		panic("原始值小于最小取样值")
	}
	result := ((src-minSpVal)/(maxSpVal-minSpVal))*(maxRangeVal-minRangeVal) + (minRangeVal)

	if result > maxRangeVal || result < minRangeVal {
		panic("计算结果超出范围")
	}
	
	return math.Round(result*100) / 100
}

println(HandlerModbusData(1024.0))`
	v, err := RunScript("main.gop", s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestRunGopFunc(t *testing.T) {
	s := `import "math"

// HandlerModbusData 
// 计算传感器值 计算公式 = （（采集值-最小取样值）/（最大取样值-最小取样值）*（最大量程-最小量程） ）  + 最小量程
func HandlerModbusData(src float64) float64 {
	var minSpVal float64 = 819.0 // 最小取样值 如 819
	var maxSpVal float64 = 4095.0 // 最大取样值 如 4095
	var minRangeVal float64 = 0.0 //最小量程 0
	var maxRangeVal float64 = 300 //最大量程 300
	if src < minSpVal {
		panic("原始值小于最小取样值")
	}
	result := ((src-minSpVal)/(maxSpVal-minSpVal))*(maxRangeVal-minRangeVal) + (minRangeVal)

	if result > maxRangeVal || result < minRangeVal {
		panic("计算结果超出范围")
	}
	
	return math.Round(result*100) / 100
}`
	v, err := RunFunc("HandlerModbusData.gop", s, "HandlerModbusData", 1024.0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func Benchmark(b *testing.B) {
	s := `import "math"

// HandlerModbusData 
// 计算传感器值 计算公式 = （（采集值-最小取样值）/（最大取样值-最小取样值）*（最大量程-最小量程） ）  + 最小量程
func HandlerModbusData(src float64) float64 {
	var minSpVal float64 = 819.0 // 最小取样值 如 819
	var maxSpVal float64 = 4095.0 // 最大取样值 如 4095
	var minRangeVal float64 = 0.0 //最小量程 0
	var maxRangeVal float64 = 300 //最大量程 300
	if src < minSpVal {
		panic("原始值小于最小取样值")
	}
	result := ((src-minSpVal)/(maxSpVal-minSpVal))*(maxRangeVal-minRangeVal) + (minRangeVal)

	if result > maxRangeVal || result < minRangeVal {
		panic("计算结果超出范围")
	}
	
	return math.Round(result*100) / 100
}`
	for i := 0; i < b.N; i++ {
		_, err := RunFunc("HandlerModbusData.gop", s, "HandlerModbusData", 1024.0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
