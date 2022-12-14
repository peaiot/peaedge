package modbus_task

import (
	"fmt"
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
)

type Regdata struct {
	DeviceId string `json:"device_id"`
	VarId    string `json:"var_id" `
	RegId    string `json:"reg_id"`
	MN       string `json:"mn" `
	Name     string `json:"name"`
	Factor   string `json:"factor"`
	Value    string `json:"value"`
}

// RegisterSaveRtdTask 定时保存实时数据到数据库
func RegisterSaveRtdTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Modbus.Error(err)
		}
	}()
	var datas = make([]Regdata, 0)
	err := app.DB().Raw(`select d.id         as device_id,
		       r.id           as reg_id,
		       v.id           as var_id,
		       d.mn           as mn,
		       v.name         as name,
		       v.data_factor  as factor,
		       r.rtd          as value
		from modbus_device d,
		     modbus_reg r,
		     modbus_var v
		where d.id == r.device_id
		  and r.var_id = v.id
		  and r.flag = 'N'
		   or r.flag == ''`).Scan(&datas).Error
	common.Must(err)

	var rtds []models.DeviceRtdData

	for _, val := range datas {
		if val.MN == "" || val.Name == "" || val.Value == "" {
			continue
		}
		rtds = append(rtds, models.DeviceRtdData{
			ID:        common.UUID(),
			MN:        val.MN,
			Name:      val.Name,
			Factor:    val.Factor,
			Value:     val.Value,
			CreatedAt: time.Now(),
		})

		// insert timeseries data
		err = app.TsDB().InsertRows([]tstorage.Row{
			{
				Metric: fmt.Sprintf("modbus_metrics_%s_%s", val.MN, val.RegId),
				DataPoint: tstorage.DataPoint{
					Value:     cast.ToFloat64(val.Value),
					Timestamp: time.Now().Unix(),
				},
			},
		})
		if err != nil {
			log.Error("add timeseries data error:", err.Error())
		}
	}

	if len(rtds) > 0 {
		err := app.DB().Create(&rtds).Error
		common.Must(err)
	}

}
