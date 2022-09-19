package mqttc

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/toughstruct/peaedge/common/log"
)

// 处理 initialize 初始化事件， 由平台下发
func (d MqttDaemon) initialize(client mqtt.Client, msg mqtt.Message) {
	var _msg InitializeMsg
	err := json.Unmarshal(msg.Payload(), &_msg)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(_msg.AccessToken)
}
