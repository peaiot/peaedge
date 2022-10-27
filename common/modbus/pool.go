package modbus

import (
	"sync"
)

type TcpTransporterPool struct {
	m     sync.Mutex
	len   int
	cPool map[string][]*TCPClientTransporter
}

func NewTcpTransporterPool(len int) *TcpTransporterPool {
	return &TcpTransporterPool{
		len:   len,
		cPool: make(map[string][]*TCPClientTransporter),
	}
}

func (pl *TcpTransporterPool) Get(address string) *TCPClientTransporter {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.cPool[address])
	if n == 0 {
		pl.cPool[address] = make([]*TCPClientTransporter, 0, pl.len)
		return pl.New(address)
	}
	x := pl.cPool[address][n-1]
	pl.cPool[address] = pl.cPool[address][0 : n-1]
	return x
}

func (pl *TcpTransporterPool) New(address string) *TCPClientTransporter {
	t := NewTCPClientTransporter(address)
	return t
}

func (pl *TcpTransporterPool) Put(address string, c *TCPClientTransporter) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.cPool[address] = append(pl.cPool[address], c)
}

func (pl *TcpTransporterPool) Shutdown() {
	for _, cs := range pl.cPool {
		for _, c := range cs {
			c.Close()
		}
	}
}

type RtuTransporterPool struct {
	m     sync.Mutex
	cPool map[string]*RTUClientTransporter
}

func NewRtuTransporterPool() *RtuTransporterPool {
	return &RtuTransporterPool{
		cPool: make(map[string]*RTUClientTransporter),
	}
}

func (pl *RtuTransporterPool) Get(address string) *RTUClientTransporter {
	pl.m.Lock()
	defer pl.m.Unlock()
	v, ok := pl.cPool[address]
	if !ok {
		return pl.New(address)
	}
	return v
}

func (pl *RtuTransporterPool) New(address string) *RTUClientTransporter {
	t := NewRTUClientTransporter(address)
	return t
}

func (pl *RtuTransporterPool) Put(address string, c *RTUClientTransporter) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.cPool[address] = c
}

func (pl *RtuTransporterPool) Shutdown() {
	for _, c := range pl.cPool {
		_ = c.Close()
	}
}
