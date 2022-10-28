package httpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/guonaihong/gout"
	_ "github.com/guonaihong/gout"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
)

type HttpClient struct {
	sync.Mutex
	Sid     string
	Debug   bool
	Url     string
	Format  string
	Header  string
	Timeout int
	Enabled bool
}

type BaseResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var clients []*HttpClient

func StartAll() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	var chls []models.HttpChannel
	err := app.DB().Where("status = 1").Find(&chls).Error
	if err != nil {
		log.Error("no http client start")
		return err
	}

	for _, chl := range chls {
		client, err := newHttpClient(chl)
		if err != nil {
			log.Errorf("restart http client -> %s error: %s", chl.Url, err)
			continue
		}
		clients = append(clients, client)
		log.Infof("restart http client -> %s success", chl.Url)
	}

	// 订阅消息
	_ = app.EvBUS().SubscribeAsync(app.EventChannelMessagePublish, func(chlType string, sid string, msg string) {
		for _, client := range clients {
			if chlType != "http" || client.Sid != sid {
				continue
			}
			r, err := client.SendBody(msg)
			if err != nil {
				log.Sched.Errorf("send http message error: %s", err)
				continue
			}
			log.Sched.Infof("send http message result: %v", r)
		}
	}, false)
	return nil
}

func newHttpClient(chl models.HttpChannel) (*HttpClient, error) {
	h := &HttpClient{
		Sid:     chl.ID,
		Url:     chl.Url,
		Format:  chl.Format,
		Header:  chl.Header,
		Timeout: chl.Timeout,
		Enabled: chl.Status == 1,
		Debug:   chl.Debug == 1,
	}
	return h, nil
}

func RestartAll() error {
	var chls []models.HttpChannel
	err := app.DB().Where("status = 1").Find(&chls).Error
	if err != nil {
		log.Error("no http client start")
		return err
	}

	clients = make([]*HttpClient, 0)
	for _, chl := range chls {
		client, err := newHttpClient(chl)
		if err != nil {
			log.Errorf("restart http client -> %s error: %s", chl.Url, err)
			continue
		}
		clients = append(clients, client)
		log.Infof("restart http client -> %s success", chl.Url)
	}
	return nil
}

func (h *HttpClient) SendBody(msg string) (*BaseResult, error) {
	resp := new(BaseResult)
	err := gout.
		POST(common.UrlJoin(h.Url)).
		SetHeader(gout.H{"Content-Type": "application/json"}).
		Debug(h.Debug).
		SetBody(msg).
		SetTimeout(time.Second * time.Duration(h.Timeout)).
		BindJSON(resp).
		Do()

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HttpClient) SendJsonData(msg models.DeviceMessage) (*BaseResult, error) {
	resp := new(BaseResult)
	err := gout.
		POST(common.UrlJoin(h.Url)).
		// SetHeader(gout.H{"accessKey": api.Apikey}).
		Debug(h.Debug).
		SetQuery(msg).
		SetTimeout(time.Second * time.Duration(h.Timeout)).
		BindJSON(resp).
		Do()

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HttpClient) SendFormData(msg models.DeviceMessage) (*BaseResult, error) {
	resp := new(BaseResult)
	err := gout.
		POST(common.UrlJoin(h.Url)).
		// SetHeader(gout.H{"accessKey": api.Apikey}).
		Debug(h.Debug).
		SetQuery(gout.H{"sign": msg.Sign}).
		SetForm(gout.H{"mn": msg.Mn}).
		SetForm(msg.Data).
		SetTimeout(time.Second * time.Duration(h.Timeout)).
		BindJSON(resp).
		Do()

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 通道消息处理
func (h *HttpClient) onChannelMessage(msg models.DeviceMessage) {
	switch h.Format {
	case "json":
		r, err := h.SendJsonData(msg)
		if err != nil {
			log.Errorf("send json message error: %s", err)
		}
		log.Infof("send json message result: %v", r)
	case "param":
		r, err := h.SendFormData(msg)
		if err != nil {
			log.Errorf("send form message error: %s", err)
		}
		log.Infof("send form message result: %v", r)
	}
}
