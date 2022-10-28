package jobs

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/jobs/datastream"
	"github.com/toughstruct/peaedge/jobs/modbus_task"
	"github.com/toughstruct/peaedge/log"
)

// Sched 计划任务管理
var Sched *cron.Cron

// Init 初始化任务计划
func Init() {
	if !app.IsInit() {
		log.Fatal("app not init")
		return
	}
	loc, _ := time.LoadLocation(app.Config().System.Location)
	Sched = cron.New(cron.WithLocation(loc))
	go modbus_task.StartModbusReadTask()

	_, _ = Sched.AddFunc(fmt.Sprintf("@every %ds", 60), func() {
		modbus_task.RegisterSaveRtdTask()
	})

	// 每分钟执行一次
	_, _ = Sched.AddFunc("* * * * *", func() {
		datastream.ProcessDatastreamTask("minute")
	})

	// 每5分钟执行一次
	_, _ = Sched.AddFunc("*/5 * * * *", func() {
		datastream.ProcessDatastreamTask("5minute")
	})

	// 每10分钟执行一次
	_, _ = Sched.AddFunc("*/10 * * * *", func() {
		datastream.ProcessDatastreamTask("10minute")
	})

	// 每小时执行一次
	_, _ = Sched.AddFunc("0 */1 * * *", func() {
		datastream.ProcessDatastreamTask("hour")
	})

	// 每天执行一次
	_, _ = Sched.AddFunc("0 0 * * *", func() {
		datastream.ProcessDatastreamTask("daily")
	})

	Sched.Start()
}
