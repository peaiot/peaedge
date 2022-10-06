package modbus

import (
	"container/ring"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"sync"

	"github.com/montanaflynn/stats"
	"github.com/shopspring/decimal"
)

const (
	// BigEndian big endian	aa bb cc dd	高尾端，高字节在低地址
	BigEndian = "BigEndian"
	// BigEndianSwap big endian byte swap	bb aa dd cc	每 2 字节使用高尾端，2 字节内部使用低尾端
	BigEndianSwap = "BigEndianSwap"
	// LittleEndian little endian	dd cc bb aa	低尾端，低字节在低地址
	LittleEndian = "LittleEndian"
	// LittleEndianSwap little endian byte swap	cc dd aa bb	每 2 字节使用低尾端，2 字节内部使用高尾端
	LittleEndianSwap = "LittleEndianSwap"
)

var DataRangeError = errors.New("数据超过量程")

var DataOverrunsError = errors.New("数据超过标准值")

var NADecimal, _ = decimal.NewFromString("999999999.9999999")

// 计算传感器值 计算公式 = （（采集值-最小取样值）/（最大取样值-最小取样值）*（最大量程-最小量程） ）  + 最小量程
// regval 采集值
// minSpVal 最小取样值 如 819
// maxSpVal 最大取样值 如 4095
// minRangeVal 最小量程
// maxRangeVal 最大量程
func ComputeResult(regval, minSpVal, maxSpVal, minRangeVal, maxRangeVal, maxAlarmVal, minAlarmVal decimal.Decimal, fixed int32, sign int) (string, error) {
	if regval.LessThan(minSpVal) {
		return "0", errors.New("采集值低于最小取样值")
	}
	result := (regval.Sub(minSpVal)).Div(maxSpVal.Sub(minSpVal)).Mul(maxRangeVal.Sub(minRangeVal)).Add(minRangeVal)

	if result.LessThan(decimal.NewFromInt(0)) {
		switch sign {
		case 1:
			result = result.Abs()
		case 100:
			result = decimal.NewFromInt(0)
		case -1:
		default:
		}
	}

	if result.GreaterThan(maxRangeVal) || result.LessThan(minRangeVal) {
		return "N/A", DataRangeError
	}

	if maxAlarmVal != NADecimal && minAlarmVal != NADecimal {
		if result.GreaterThan(maxAlarmVal) || result.LessThan(minAlarmVal) {
			return result.StringFixed(fixed), DataOverrunsError
		}
	}
	return result.StringFixed(fixed), nil
}

func DoComputeResult(regval float64, minSpVal string, maxSpVal string, minRangeVal, maxRangeVal, maxAlarmVal, minAlarmVal string, fixed int32, sign int) (string, error) {
	_regval := decimal.NewFromFloat(regval)

	_minSpVal, err := decimal.NewFromString(minSpVal)
	if err != nil {
		return "", err
	}

	_maxSpVal, err := decimal.NewFromString(maxSpVal)
	if err != nil {
		return "", err
	}

	_minRangeVal, err := decimal.NewFromString(minRangeVal)
	if err != nil {
		return "", err
	}
	_maxRangeVal, err := decimal.NewFromString(maxRangeVal)
	if err != nil {
		return "", err
	}

	_maxAlarmVal, err := decimal.NewFromString(maxAlarmVal)
	if err != nil {
		_maxAlarmVal = NADecimal
	}
	_minAlarmVal, err := decimal.NewFromString(minAlarmVal)
	if err != nil {
		_minAlarmVal = NADecimal
	}

	return ComputeResult(_regval, _minSpVal, _maxSpVal, _minRangeVal, _maxRangeVal, _maxAlarmVal, _minAlarmVal, fixed, sign)
}

// 根据数字系数计算
func DoComputeDxyResult(regval float64, dx string, dy string, minRangeVal, maxRangeVal, maxAlarmVal, minAlarmVal string, fixed int32, sign int) (string, error) {
	_regval := decimal.NewFromFloat(regval)
	_x, err := decimal.NewFromString(dx)
	if err != nil {
		return "", err
	}
	_y, err := decimal.NewFromString(dy)
	if err != nil {
		return "", err
	}

	_minRangeVal, err := decimal.NewFromString(minRangeVal)
	if err != nil {
		return "", err
	}
	_maxRangeVal, err := decimal.NewFromString(maxRangeVal)
	if err != nil {
		return "", err
	}

	_maxAlarmVal, err := decimal.NewFromString(maxAlarmVal)
	if err != nil {
		_maxAlarmVal = NADecimal
	}
	_minAlarmVal, err := decimal.NewFromString(minAlarmVal)
	if err != nil {
		_minAlarmVal = NADecimal
	}

	result := _regval.Mul(_x).Add(_y)
	if result.LessThan(decimal.NewFromInt(0)) {
		switch sign {
		case 1:
			result = result.Abs()
		case 100:
			result = decimal.NewFromInt(0)
		case -1:
		default:
		}
	}

	if result.GreaterThan(_maxRangeVal) || result.LessThan(_minRangeVal) {
		return "N/A", DataRangeError
	}

	if _maxAlarmVal != NADecimal && _minAlarmVal != NADecimal {
		if result.GreaterThan(_maxAlarmVal) || result.LessThan(_minAlarmVal) {
			return result.StringFixed(fixed), DataOverrunsError
		}
	}

	return result.StringFixed(fixed), nil
}

func GetFloatValue(p1, p2 uint16) float64 {
	var intSign, intSignRest, intExponent, intExponentRest int
	var faResult float64
	var faDigit float64
	intSign = int(p1 / 32768)
	intSignRest = int(p2 % 32768)
	intExponent = intSignRest / 128
	intExponentRest = intSignRest % 128
	faDigit = float64(intExponentRest*65536+int(p1)) / 8388608
	v, _ := stats.Round(faDigit, 9)
	faResult = math.Pow(-1, float64(intSign)) * math.Pow(2, float64(intExponent)-127) * (v + 1)
	return faResult
}

// float getFloat(quint16 value1, quint16 value2)
// {
//     float fTemp;
//     uint *pTemp=(uint *)&fTemp;
//     unsigned int chTemp[4];//a,b,c,d
//     chTemp[0]=value1&0xff;
//     chTemp[1]=(value1>>8)&0xff;
//     chTemp[2]=value2&0xff;
//     chTemp[3]=(value2>>8)&0xff;
//     //这是ABCD
//     *pTemp=((chTemp[1]<<24)&0xff000000)|((chTemp[0]<<16)&0xff0000)|((chTemp[3]<<8)&0xff00)|(chTemp[2]&0xff);
//
//     //这是CDAB
//     //*pTemp=((chTemp[3]<<24)&0xff000000)|((chTemp[2]<<16)&0xff0000)|((chTemp[1]<<8)&0xff00)|(chTemp[0]&0xff);
//
//     //这是BADC
//     //*pTemp=((chTemp[0]<<24)&0xff000000)|((chTemp[1]<<16)&0xff0000)|((chTemp[2]<<8)&0xff00)|(chTemp[3]&0xff);
//
//     //这是DCBA
//     //*pTemp=((chTemp[2]<<24)&0xff000000)|((chTemp[3]<<16)&0xff0000)|((chTemp[0]<<8)&0xff00)|(chTemp[1]&0xff);
//     return fTemp;
// }

func GetFloatValue2(p1, p2 uint16) float32 {
	var chTemp []uint32 = make([]uint32, 4)
	chTemp[0] = uint32(p1 & 0xff)
	chTemp[1] = uint32((p1 >> 8) & 0xff)
	chTemp[2] = uint32(p2 & 0xff)
	chTemp[3] = uint32((p2 >> 8) & 0xff)

	var pTemp = ((chTemp[1] << 24) & 0xff000000) | ((chTemp[0] << 16) & 0xff0000) | ((chTemp[3] << 8) & 0xff00) | (chTemp[2] & 0xff)
	return math.Float32frombits(pTemp)
}

// def convert_registers_to_float(registers):
//     """
//     Convert two 16 Bit Registers to 32 Bit real value - Used to receive float values from Modbus (Modbus Registers are 16 Bit long)
//     registers: 16 Bit Registers
//     return: 32 bit value real
//     """
//     b = bytearray(4)
//     b [0] = registers[0] & 0xff
//     b [1] = (registers[0] & 0xff00)>>8
//     b [2] = (registers[1] & 0xff)
//     b [3] = (registers[1] & 0xff00)>>8
//     returnValue = struct.unpack('<f', b)            #little Endian
//     return returnValue
//

// GetFloat32Value
/*
字节序	内存低地址 -> 内存高地址	备注
big endian	aa bb cc dd	高尾端，高字节在低地址
big endian byte swap	bb aa dd cc	每 2 字节使用高尾端，2 字节内部使用低尾端
little endian	dd cc bb aa	低尾端，低字节在低地址
little endian byte swap	cc dd aa bb	每 2 字节使用低尾端，2 字节内部使用高尾端
*/
func GetFloat32Value(p1, p2 uint16, format string) float32 {
	var A, B, C, D byte
	switch format {
	case BigEndian: // ABCD
		A, B, C, D = byte(p1&0xff), byte((p1&0xff00)>>8), byte(p2&0xff), byte((p2&0xff00)>>8)
	case BigEndianSwap: // BADC
		B, A, D, C = byte(p1&0xff), byte((p1&0xff00)>>8), byte(p2&0xff), byte((p2&0xff00)>>8)
	case LittleEndian: // DCBA
		D, C, B, A = byte(p1&0xff), byte((p1&0xff00)>>8), byte(p2&0xff), byte((p2&0xff00)>>8)
	case LittleEndianSwap: // CDAB
		C, D, A, B = byte(p1&0xff), byte((p1&0xff00)>>8), byte(p2&0xff), byte((p2&0xff00)>>8)
	default:
		A, B, C, D = byte(p1&0xff), byte((p1&0xff00)>>8), byte(p2&0xff), byte((p2&0xff00)>>8)
	}
	return math.Float32frombits(binary.LittleEndian.Uint32([]byte{A, B, C, D}))
}

// GetFloat64Value
/*
字节序	内存低地址 -> 内存高地址	备注
big endian	aa bb cc dd	高尾端，高字节在低地址
big endian byte swap	bb aa dd cc	每 2 字节使用高尾端，2 字节内部使用低尾端
little endian	dd cc bb aa	低尾端，低字节在低地址
little endian byte swap	cc dd aa bb	每 2 字节使用低尾端，2 字节内部使用高尾端
*/
func GetFloat64Value(p1, p2, p3, p4 uint16, format string) float64 {
	V1 := byte(p1 & 0xff)
	V2 := byte((p1 & 0xff00) >> 8)
	V3 := byte(p2 & 0xff)
	V4 := byte((p2 & 0xff00) >> 8)
	V5 := byte(p3 & 0xff)
	V6 := byte((p3 & 0xff00) >> 8)
	V7 := byte(p4 & 0xff)
	V8 := byte((p4 & 0xff00) >> 8)
	var A, B, C, D, E, F, G, H byte
	switch format {
	case BigEndian: // ABCDEFGH:
		A, B, C, D, E, F, G, H = V1, V2, V3, V4, V5, V6, V7, V8
	case BigEndianSwap: // BADCFEHG:
		B, A, D, C, F, E, H, G = V1, V2, V3, V4, V5, V6, V7, V8
	case LittleEndian: // HGFEDCBA
		H, G, F, E, D, C, B, A = V1, V2, V3, V4, V5, V6, V7, V8
	case LittleEndianSwap: // GHEFCDAB
		G, H, E, F, C, D, A, B = V1, V2, V3, V4, V5, V6, V7, V8
	default:
		A, B, C, D, E, F, G, H = V1, V2, V3, V4, V5, V6, V7, V8
	}
	return math.Float64frombits(binary.LittleEndian.Uint64([]byte{A, B, C, D, E, F, G, H}))
}

// func GetFloatValue3(p1, p2 uint16) float64 {
// 	b := make([]byte, 4)
// 	b[0] = byte(p1 & 0xff)
// 	b[1] = byte((p1 & 0xff00) >> 8)
// 	b[2] = byte(p2 & 0xff)
// 	b[3] = byte((p2 & 0xff00) >> 8)
// 	return float64(math.Float32frombits(binary.LittleEndian.Uint32(b)))
// }

type TcpVlidate struct {
	lock   *sync.Mutex
	valmap map[int64]*ring.Ring
}

func NewTcpVlidate() *TcpVlidate {
	return &TcpVlidate{valmap: make(map[int64]*ring.Ring, 0), lock: &sync.Mutex{}}
}

func (tv TcpVlidate) SetLast(key int64, val float64) {
	tv.lock.Lock()
	_, ok := tv.valmap[key]
	if !ok {
		tv.valmap[key] = ring.New(3)
	}
	tv.valmap[key].Value = val
	tv.valmap[key] = tv.valmap[key].Next()
	tv.lock.Unlock()
}

func (tv TcpVlidate) Sum(key int64) float64 {
	tv.lock.Lock()
	_, ok := tv.valmap[key]
	if !ok {
		tv.lock.Unlock()
		return 0
	}
	var r float64 = 0
	tv.valmap[key].Do(func(p interface{}) {
		r = r + p.(float64)
	})
	tv.lock.Unlock()
	return r
}

func (tv TcpVlidate) Len(key int64) int {
	tv.lock.Lock()
	_, ok := tv.valmap[key]
	if !ok {
		tv.lock.Unlock()
		return 0
	}
	var l = 0
	tv.valmap[key].Do(func(p interface{}) {
		v := reflect.ValueOf(p)
		if v.Kind() == reflect.Float64 {
			l = l + 1
		}
	})
	tv.lock.Unlock()
	return l
}
