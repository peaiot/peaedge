package golua

import (
	"errors"
	"go/types"

	"github.com/montanaflynn/stats"
	"github.com/spf13/cast"
	"github.com/toughstruct/peaedge/golua/appx"
	"github.com/toughstruct/peaedge/golua/json"
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

func GetGoValue(arg lua.LValue) any {
	switch arg.(type) {
	case lua.LNumber:
		return float64(arg.(lua.LNumber))
	case lua.LString:
		return arg.String()
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(arg.(lua.LBool))
	}
	return arg.String()
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
func HandlerMixedRegister(script string, regs map[string]string, vars map[string]string) (ret float64, err error) {
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
		rtab.RawSetString(rk, lua.LString(rv))
	}
	vtab := L.NewTable()
	for vk, vv := range vars {
		vtab.RawSetString(vk, lua.LString(vv))
	}
	L.Push(rtab)
	L.Push(vtab)
	L.Call(2, 1)
	r := L.Get(L.GetTop())
	ret, err = stats.Round(cast.ToFloat64(r.String()), 2)
	return
}

// HandlerDataStream 数据流处理，最终输出字符串
func HandlerDataStream(script string, arg any) (ret string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = err2.(error)
		}
	}()
	L := lua.NewState()
	defer L.Close()
	appx.LuaPreload(L)
	json.Preload(L)
	if err = L.DoString(script); err != nil {
		return
	}
	fn := L.GetGlobal("HandlerDataStream").(*lua.LFunction)
	L.Push(fn)
	L.Push(GetLuaValue(arg))
	L.Call(1, 1)
	r := L.Get(L.GetTop())
	_, ok := r.(lua.LString)
	if !ok {
		return "", errors.New("return value is not string")
	}
	ret = r.String()
	return
}
