package metrics

import (
	"sync"
	"sync/atomic"
)

type Metrics struct {
	Name  string
	Value int64
}

func (m Metrics) Incr(value int64) {
	atomic.AddInt64(&m.Value, value)
}

type ValueData struct {
	sync.Mutex
	values map[string]Metrics
}

func NewValueData() *ValueData {
	return &ValueData{values: make(map[string]Metrics)}
}

func (md *ValueData) Incr(name string, value int64) {
	md.Lock()
	defer md.Unlock()
	if _, ok := md.values[name]; ok {
		md.values[name].Incr(value)
	} else {
		md.values[name] = Metrics{Name: name, Value: value}
	}
}

func (md *ValueData) Get(name string) int64 {
	md.Lock()
	defer md.Unlock()
	if v, ok := md.values[name]; ok {
		return v.Value
	}
	return 0
}
