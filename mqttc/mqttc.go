package mqttc

import (
	"fmt"
	slog "log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/common/sysid"
)

var daemon *MqttDaemon

func StartDaemon() error {
	daemon = newMqttDaemon()
	return daemon.Start()
}

type MqttDaemon struct {
	BrokerSerer string
	ClientId    string
	Mqttc       mqtt.Client
}

type InitializeMsg struct {
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

func newMqttDaemon() *MqttDaemon {
	d := &MqttDaemon{}
	d.ClientId = sysid.GetSystemSid()
	d.BrokerSerer = app.Config.Mqtt.Broker
	return d
}

func (d MqttDaemon) Publish(topic string, qos int, payload interface{}) {
	t := d.Mqttc.Publish(topic, byte(qos), false, payload)
	if t.Wait() && t.Error() != nil {
		log.Errorf("publish %s error: %s\n", topic, t.Error())
	}
}

func (d MqttDaemon) onConnect(client mqtt.Client) {
	log.Infof("connect %v", client)
	t := client.Subscribe(PeaEdgeInitialize, 2, d.initialize)
	if t.Wait() && t.Error() != nil {
		panic(fmt.Errorf("subscribe %s error: %s\n", PeaEdgeInitialize, t.Error()))
	}
	go d.bootstrap()
}

// Start 启动守护进程
func (d MqttDaemon) Start() error {
	mqtt.DEBUG = slog.New(log.Stdlog{}, "Mqtt", 0)
	mqtt.ERROR = slog.New(log.Stderr{}, "Mqtt", 0)
	opts := mqtt.NewClientOptions().
		AddBroker(d.BrokerSerer).
		SetClientID(d.ClientId).
		SetUsername(app.Config.Mqtt.Username).
		SetPassword(app.Config.Mqtt.Password)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.OnConnect = d.onConnect

	d.Mqttc = mqtt.NewClient(opts)
	if token := d.Mqttc.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
