package monitor

import (
	"os"
	"time"

	"github.com/nakabonne/tstorage"
	_ "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/log"
)

func ProcessMonitorTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Sched.Error(err)
		}
	}()

	timestamp := time.Now().Unix()

	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return
	}

	cpuuse, _ := p.CPUPercent()
	if err != nil {
		cpuuse = 0
	}

	err = app.TsDB().InsertRows([]tstorage.Row{
		{
			Metric: "peaedge_cpuuse",
			DataPoint: tstorage.DataPoint{
				Value:     cpuuse,
				Timestamp: timestamp,
			},
		},
	})
	if err != nil {
		log.Error("add timeseries data error:", err.Error())
	}

	meminfo, err := p.MemoryInfo()
	if err != nil {
		return
	}
	memuse := meminfo.RSS / 1024 / 1024

	err = app.TsDB().InsertRows([]tstorage.Row{
		{
			Metric: "peaedge_memuse",
			DataPoint: tstorage.DataPoint{
				Value:     float64(memuse),
				Timestamp: timestamp,
			},
		},
	})
	if err != nil {
		log.Error("add timeseries data error:", err.Error())
	}

}
