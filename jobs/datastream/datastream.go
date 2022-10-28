package datastream

import (
	"strings"

	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/golua"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
)

func ProcessDatastreamTask(schedPolicy string) {
	defer func() {
		if err := recover(); err != nil {
			log.Sched.Error(err)
		}
	}()

	var datastreams []models.DataStream
	err := app.DB().Where("sched_policy = ?", schedPolicy).Find(&datastreams).Error
	if err != nil {
		log.Sched.Errorf("读取数据流列表失败... %s", err.Error())
		return
	}

	for _, ds := range datastreams {
		var script models.DataScript
		err = app.DB().
			Where("id = ? and func_name = ?", ds.ScriptId, app.LuaFuncHandlerDataStream).
			First(&script).Error
		if err != nil {
			log.Sched.Errorf("读取数据流脚本失败... %s", err.Error())
			continue
		}
		result, err := golua.HandlerDataStream(script.Content, ds.MN)
		if err != nil {
			log.Sched.Errorf("数据流脚本执行失败... %s", err.Error())
			continue
		}

		for _, cid := range strings.Split(ds.MqttChids, ",") {
			app.PubChannelMessage("mqtt", cid, result)
		}

		for _, cid := range strings.Split(ds.HttpChids, ",") {
			app.PubChannelMessage("http", cid, result)
		}

		for _, cid := range strings.Split(ds.TcpChids, ",") {
			app.PubChannelMessage("tcp", cid, result)
		}

	}

}
