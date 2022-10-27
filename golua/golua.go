package golua

import (
	"go/types"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cast"
	lua "github.com/yuin/gopher-lua"
)

func GetLuaValue(arg any) lua.LValue {
	switch arg.(type) {
	case int, int16, int32, uint, uint32, uint16, uint64, int64, int8, uint8, float64, float32:
		return lua.LNumber(arg.(int))
	case string:
		return lua.LString(arg.(string))
	case types.Nil:
		return lua.LNil
	case bool:
		return lua.LBool(arg.(bool))
	}
	return lua.LNil
}

// HandlerModbusRtd modbus 实时数据计算处理
// arg 输入数值
// ret 计算结果 float64
func HandlerModbusRtd(script string, arg any) (ret float64, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = err2.(error)
		}
	}()
	L := lua.NewState()
	defer L.Close()
	if err = L.DoString(script); err != nil {
		return
	}
	fn := L.GetGlobal("HandlerModbusRtd").(*lua.LFunction)
	L.Push(fn)
	L.Push(GetLuaValue(arg))
	L.Call(1, 1)
	r := L.Get(L.GetTop())
	ret, err = stats.Round(cast.ToFloat64(r.String()), 2)
	return
}

// HandlerMixedRegister 混合寄存器计算
// regs 寄存器值映射表
// vtab 变量值映射表
// ret 计算结果 float64
func HandlerMixedRegister(script string, regs map[string]float64, vars map[string]float64) (ret float64, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = err2.(error)
		}
	}()
	L := lua.NewState()
	defer L.Close()
	if err = L.DoString(script); err != nil {
		return
	}
	fn := L.GetGlobal("HandlerMixedRegister").(*lua.LFunction)
	L.Push(fn)
	rtab := L.NewTable()
	for rk, rv := range regs {
		rtab.RawSetString(rk, lua.LNumber(rv))
	}
	vtab := L.NewTable()
	for vk, vv := range vars {
		vtab.RawSetString(vk, lua.LNumber(vv))
	}
	L.Push(rtab)
	L.Push(vtab)
	L.Call(2, 1)
	r := L.Get(L.GetTop())
	ret, err = stats.Round(cast.ToFloat64(r.String()), 2)
	return
}
