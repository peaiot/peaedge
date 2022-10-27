package app

import (
	"time"

	"github.com/toughstruct/peaedge/models"
)

const (
	EventMqttcBootstrap        = "mqttc_bootstrap"
	EventChannelMessagePublish = "channel_message_publish"
)

func setupSubscribers() {
	_ = evBus.SubscribeAsync(EventMqttcBootstrap, func(sid string) {
		DB().Model(&models.MqttChannel{}).Where("id = ?", sid).Update("last_boot", time.Now())
	}, false)

}

// PubChannelMessage 发布通道消息
func PubChannelMessage(msgType string, msg any) {
	evBus.Publish(EventChannelMessagePublish, msgType, msg)
}
