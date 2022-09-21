package jobs

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common/timeutil"
	"github.com/toughstruct/peaedge/jobs/modbus_task"
)

// 计划任务管理

var Sched *cron.Cron
var nameMapLock sync.Mutex
var nameMap = make(map[cron.EntryID]string, 0)
var Names = make([]string, 0)

type namedJob struct {
	Id  cron.EntryID
	err error
}

func newNamedJob(id cron.EntryID, err error) *namedJob {
	return &namedJob{Id: id, err: err}
}

func addNamedJob(name string, ljob *namedJob) {
	nameMapLock.Lock()
	defer nameMapLock.Unlock()
	nameMap[ljob.Id] = name
	Names = append(Names, name)
}

// Init 初始化任务计划
func Init() {
	loc, _ := time.LoadLocation(app.Config.System.Location)
	go modbus_task.StartModbusReadTask()
	Sched = cron.New(cron.WithLocation(loc))
	Sched.Start()
}

// QueryJobEntry 提供 WEB 查询
func QueryJobEntry() []map[string]interface{} {
	nameMapLock.Lock()
	defer nameMapLock.Unlock()
	result := make([]map[string]interface{}, 0)
	for _, ent := range Sched.Entries() {
		name, ok := nameMap[ent.ID]
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":        ent.ID,
			"remark":    name,
			"last_exec": ent.Prev.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
			"next_exec": ent.Next.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
		})
	}
	return result
}

var wg sync.WaitGroup

// ExecJob 立即运行任务
func ExecJob(id int) {
	ent := Sched.Entry(cron.EntryID(id))
	if ent.Valid() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ent.Job.Run()
		}()
	}
}
