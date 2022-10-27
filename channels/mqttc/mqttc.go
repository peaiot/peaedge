package mqttc

import (
	"fmt"
	slog "log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
)

var daemons []*MqttDaemon

func StartAll() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	mqtt.DEBUG = slog.New(log.Mqttc, "Mqtt", 0)
	mqtt.ERROR = slog.New(log.Mqttc, "Mqtt", 0)
	var chls []models.MqttChannel
	err := app.DB().Where("status = 'enabled'").Find(&chls).Error
	if err != nil {
		log.Error("no mqtt client start")
		return nil
	}
	for _, chl := range chls {
		daemon, err := newMqttDaemon(chl)
		if err != nil {
			log.Errorf("start mqtt client %s error: %s", chl.ClientId, err)
			continue
		}
		daemons = append(daemons, daemon)
		log.Infof("start mqtt client %s success", chl.ClientId)
	}
	return nil
}

func RestartAll() error {
	var chls []models.MqttChannel
	err := app.DB().Where("status = 1").Find(&chls).Error
	if err != nil {
		log.Error("no mqtt client start")
		return nil
	}

	for _, daemon := range daemons {
		daemon.Disconnect(0)
	}
	daemons = make([]*MqttDaemon, 0)
	for _, chl := range chls {
		daemon, err := newMqttDaemon(chl)
		if err != nil {
			log.Errorf("restart mqtt client %s error: %s", chl.ClientId, err)
			continue
		}
		daemons = append(daemons, daemon)
		log.Infof("restart mqtt client %s success", chl.ClientId)
	}
	return nil
}

type MqttDaemon struct {
	Sid             string
	BrokerSerer     string
	ClientId        string
	Username        string
	Password        string
	KeepAlive       int
	PingTimeout     int
	RetryInterval   int
	ClearSession    bool
	ProtocolVersion int
	SubTopic        string
	PubTopic        string
	Retained        bool
	Qos             int
	Will            string
	Enabled         bool
	Mqttc           mqtt.Client
}

func (d MqttDaemon) Publish(qos int, payload interface{}) error {
	if !d.Enabled {
		return nil
	}
	t := d.Mqttc.Publish(d.PubTopic, byte(qos), d.Retained, payload)
	if t.Wait() && t.Error() != nil {
		log.Mqttc.Errorf("publish %s error: %s\n", d.PubTopic, t.Error())
		return t.Error()
	}
	return nil
}

func (d MqttDaemon) onConnect(client mqtt.Client) {
	payload, _ := newBootstrapMessage(d.ClientId).Encode()
	t := client.Publish(d.PubTopic, byte(1), d.Retained, payload)
	if t.Wait() && t.Error() != nil {
		log.Mqttc.Errorf("onConnect publish %s error: %s\n", d.PubTopic, t.Error())
	}
	if token := client.Subscribe(d.SubTopic, 1, d.onMessage); token.Wait() && token.Error() != nil {
		log.Mqttc.Errorf("onConnect subscribe %s error: %s\n", d.SubTopic, t.Error())
	}

	err := app.EvBUS().SubscribeAsync(app.EventChannelMessagePublish, d.onMessage, false)
	if err != nil {
		log.Mqttc.Errorf("onConnect subscribe %s error: %s\n", app.EventChannelMessagePublish, err)
	}
}

func (d MqttDaemon) onMessage(c mqtt.Client, msg mqtt.Message) {
	// payload := msg.Payload()
	// log.Mqttc.Infof("onMessage %s %s\n", msg.Topic(), payload)
}

func (d MqttDaemon) Disconnect(quiesce uint) {
	_ = app.EvBUS().Unsubscribe(app.EventChannelMessagePublish, d.onMessage)
	d.Mqttc.Disconnect(quiesce)
}

// 通道消息处理
func (d MqttDaemon) onChannelMessage(msg models.DeviceMessage) {
	message := NewMessage[models.DeviceMessage]("dataMessage", msg)
	bs, err := message.Encode()
	if err != nil {
		log.Mqttc.Errorf("onChannelMessage encode error: %s\n", err)
		return
	}
	err = d.Publish(d.Qos, bs)
	if err != nil {
		log.Mqttc.Errorf("onChannelMessage mqtt publish error: %s\n", err)
		return
	}
}

// Start 启动守护进程
func newMqttDaemon(mc models.MqttChannel) (*MqttDaemon, error) {
	d := &MqttDaemon{
		Sid:             mc.ID,
		BrokerSerer:     mc.Server,
		ClientId:        mc.ClientId,
		Username:        mc.Username,
		Password:        mc.Password,
		KeepAlive:       mc.KeepAlive,
		PingTimeout:     mc.PingTimeout,
		RetryInterval:   mc.RetryInterval,
		ClearSession:    mc.ClearSession == 1,
		ProtocolVersion: mc.ProtocolVersion,
		SubTopic:        mc.SubTopic,
		PubTopic:        mc.PubTopic,
		Retained:        mc.Retained == 1,
		Qos:             mc.Qos,
		Will:            mc.Will,
		Enabled:         mc.Status == 1,
		Mqttc:           nil,
	}
	opts := mqtt.NewClientOptions().
		AddBroker(d.BrokerSerer).
		SetClientID(d.ClientId).
		SetUsername(d.Username).
		SetPassword(d.Password)
	opts.SetKeepAlive(time.Duration(d.KeepAlive) * time.Second)
	opts.SetPingTimeout(time.Duration(d.PingTimeout) * time.Second)
	opts.SetConnectRetryInterval(time.Duration(d.RetryInterval) * time.Second)
	opts.SetCleanSession(d.ClearSession)
	opts.SetProtocolVersion(uint(d.ProtocolVersion))
	opts.SetWill(d.PubTopic, d.Will, byte(d.Qos), d.Retained)
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = d.onConnect
	d.Mqttc = mqtt.NewClient(opts)
	if token := d.Mqttc.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return d, nil
}

func GetMqttcDaemon(sid string) *MqttDaemon {
	for _, daemon := range daemons {
		if daemon.Sid == sid {
			return daemon
		}
	}
	return nil
}
