package tcpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
)

var clients []*TcpClient

type TcpClient struct {
	sync.Mutex
	Sid         string
	Server      string
	Port        int
	ChannelType string
	Timeout     int
	Enabled     bool
	Conn        *net.TCPConn
	rw          *bufio.ReadWriter
	isStop      bool
}

func StartAll() error {
	if !app.IsInit() {
		return fmt.Errorf("app not init")
	}
	var chls []models.TcpChannel
	err := app.DB().Where("status = 1").Find(&chls).Error
	if err != nil {
		log.Error("no tcp client start")
		return err
	}

	for _, chl := range chls {
		client, err := newTcpClient(chl)
		if err != nil {
			log.Errorf("restart tcp client -> %s error: %s", chl.Server, err)
			continue
		}
		clients = append(clients, client)
		log.Infof("restart tcp client -> %s success", chl.Server)
	}
	return nil
}

func RestartAll() error {
	var chls []models.TcpChannel
	err := app.DB().Where("status = 1").Find(&chls).Error
	if err != nil {
		log.Error("no mqtt client start")
		return err
	}
	for _, client := range clients {
		client.Stop()
	}
	clients = make([]*TcpClient, 0)
	for _, chl := range chls {
		client, err := newTcpClient(chl)
		if err != nil {
			log.Errorf("start tcp client -> %s error: %s", chl.Server, err)
			continue
		}
		clients = append(clients, client)
		log.Infof("start tcp client -> %s success", chl.Server)
	}
	return nil
}

func newTcpClient(chl models.TcpChannel) (*TcpClient, error) {
	t := &TcpClient{
		Sid:         chl.ID,
		Server:      chl.Server,
		Port:        chl.Port,
		ChannelType: chl.ChannelType,
		Timeout:     chl.Timeout,
		Enabled:     chl.Status == 1,
	}
	go t.checkConn()
	return t, nil
}

// Connect 链接服务器
func (c *TcpClient) Connect() error {
	c.Lock()
	defer c.Unlock()
	addr := c.Server
	tcpAdd, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	c.Conn, err = net.DialTCP("tcp", nil, tcpAdd)
	if err != nil {
		return err
	}
	c.rw = bufio.NewReadWriter(bufio.NewReader(c.Conn), bufio.NewWriter(c.Conn))
	go c.onMessageRectived()
	return nil
}

func (c *TcpClient) DisConnect() {
	c.Lock()
	defer c.Unlock()
	_ = c.Conn.Close()
	c.Conn = nil
	c.rw = nil
}

// Stop 关闭停止客户端
func (c *TcpClient) Stop() {
	c.Lock()
	defer c.Unlock()
	_ = c.Conn.Close()
	c.Conn = nil
	c.rw = nil
	c.isStop = true
}

// 通道消息处理
func (c *TcpClient) onChannelMessage(msg models.DeviceMessage) {
	bs, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		log.Errorf("json marshal error: %s", err)
		return
	}
	err = c.SendMessage(bs)
	if err != nil {
		log.Errorf("send message error: %s", err)
	}
}

// 处理接收消息
func (c *TcpClient) onMessageRectived() {
	for {
		if c.isStop || c.Conn == nil || c.rw == nil {
			log.Errorf("server connect is nil %s", c.Server)
			return
		}
		msg, err := c.rw.ReadString('\n') // 读取直到输入中第一次发生 ‘\n’
		fmt.Println(msg)
		switch err {
		case io.EOF, net.ErrClosed, net.ErrWriteToConnected:
			log.Errorf("server %s connect error %s", c.Server, err.Error())
			c.DisConnect()
			return
		}
	}
}

// SendMessage 发送消息
func (c *TcpClient) SendMessage(message []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.rw == nil {
		return fmt.Errorf("server %s connect rw is nil", c.Server)
	}
	count, err := c.rw.Write(message)
	if err != nil {
		return err
	}
	err = c.rw.Flush()
	if err != nil {
		return err
	}
	log.Infof("send data len %d", count)
	return nil
}

// 检查链接
func (c *TcpClient) checkConn() {
	timer := time.NewTicker(time.Second * 30)
	for {
		if c.isStop {
			return
		}
		if c.Conn == nil {
			err := c.Connect()
			if err != nil {
				log.Errorf("connect server %s failure %s ", c.Server, err.Error())
			} else {
				log.Infof("connect server %s success", c.Server)
			}
		}
		// 发送心跳
		if c.rw != nil {
			err := c.SendMessage([]byte{0x00})
			switch err {
			case io.EOF, net.ErrClosed, net.ErrWriteToConnected:
				log.Errorf("server %s connect error %s", c.Server, err.Error())
				c.DisConnect()
				continue
			}
		}
		<-timer.C
	}
}
