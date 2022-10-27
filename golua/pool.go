package golua

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type LStatePool struct {
	m     sync.Mutex
	saved []*lua.LState
}

func NewLStatePool(length int) *LStatePool {
	return &LStatePool{saved: make([]*lua.LState, 0, length)}
}

func (pl *LStatePool) Get() *lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *LStatePool) New() *lua.LState {
	L := lua.NewState()
	// setting the L up here.
	// load scripts, set global variables, share channels, etc...
	return L
}

func (pl *LStatePool) Put(L *lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *LStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}

// Global LState pool
// var luaPool = &LStatePool{
//     saved: make([]*lua.LState, 0, 4),
// }
