package jobs

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"github.com/toughstruct/peaedge/app"
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

	_, _ = Sched.AddFunc(fmt.Sprintf("@every %s", app.Config().Data.RtdSave), func() {
		modbus_task.RegisterSaveRtdTask()
	})

	Sched.Start()
}
