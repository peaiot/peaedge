package appx

import (
	"github.com/toughstruct/peaedge/app"
	lua "github.com/yuin/gopher-lua"
)

// local app = require("app")

func LuaPreload(L *lua.LState) {
	L.PreloadModule("appx", Loader)
}

func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, luaapi)
	L.Push(t)
	return 1
}

var luaapi = map[string]lua.LGFunction{
	"getRegisters": getRegisters,
}

func getRegisters(L *lua.LState) int {
	mn := L.CheckString(1)
	var regs []struct {
		Factor string
		Rtd    string
	}

	err := app.DB().
		Raw(`select b.data_factor as factor,
				       a.rtd
				from modbus_reg a,
				     modbus_var b,
				     modbus_device c
				where a.var_id = b.id
				  and a.device_id = c.id
				  and c.mn =  ?`, mn).
		Scan(&regs).Error
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	for i, reg := range regs {
		subtab := L.NewTable()
		subtab.RawSetString("factor", lua.LString(reg.Factor))
		subtab.RawSetString("value", lua.LString(reg.Rtd))
		result.RawSetInt(i+1, subtab)
	}

	L.Push(result)
	return 1
}
