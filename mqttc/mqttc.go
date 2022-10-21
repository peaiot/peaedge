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

func StartDaemon() error {
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
		daemon, err := Start(chl)
		if err != nil {
			log.Errorf("start mqtt client %s error: %s", chl.ClientId, err)
			continue
		}
		daemons = append(daemons, daemon)
		log.Infof("start mqtt client %s success", chl.ClientId)
	}
	return nil
}

type MqttDaemon struct {
	BrokerSerer     string
	ClientId        string
	Username        string
	Password        string
	KeepAlive       int
	PingTimeout     int
	RetryInterval   int
	CleanSession    bool
	ProtocolVersion int
	SubTopic        string
	PubTopic        string
	Retained        bool
	Mqttc           mqtt.Client
}

func (d MqttDaemon) Publish(qos int, payload interface{}) {
	t := d.Mqttc.Publish(d.PubTopic, byte(qos), d.Retained, payload)
	if t.Wait() && t.Error() != nil {
		log.Mqttc.Errorf("publish %s error: %s\n", d.PubTopic, t.Error())
	}
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
}

func (d MqttDaemon) onMessage(c mqtt.Client, msg mqtt.Message) {
	// payload := msg.Payload()
}

// Start 启动守护进程
func Start(mc models.MqttChannel) (*MqttDaemon, error) {
	d := &MqttDaemon{
		BrokerSerer:     mc.Server,
		ClientId:        mc.ClientId,
		Username:        mc.Username,
		Password:        mc.Password,
		KeepAlive:       mc.KeepAlive,
		PingTimeout:     mc.PingTimeout,
		RetryInterval:   mc.RetryInterval,
		CleanSession:    mc.CleanSession == 1,
		ProtocolVersion: mc.ProtocolVersion,
		SubTopic:        mc.SubTopic,
		PubTopic:        mc.PubTopic,
		Retained:        mc.Retained == 1,
	}
	opts := mqtt.NewClientOptions().
		AddBroker(d.BrokerSerer).
		SetClientID(d.ClientId).
		SetUsername(d.Username).
		SetPassword(d.Password)
	opts.SetKeepAlive(time.Duration(d.KeepAlive) * time.Second)
	opts.SetPingTimeout(time.Duration(d.PingTimeout) * time.Second)
	opts.SetConnectRetryInterval(time.Duration(d.RetryInterval) * time.Second)
	opts.SetCleanSession(d.CleanSession)
	opts.SetProtocolVersion(uint(d.ProtocolVersion))
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = d.onConnect
	d.Mqttc = mqtt.NewClient(opts)
	if token := d.Mqttc.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return d, nil
}
