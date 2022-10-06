package modbus

import (
	"container/ring"
	"fmt"
	"testing"
)

func TestDoComputeResult(t *testing.T) {
	val, err := DoComputeResult(1024, "819", "4095", "0", "300", "", "", 2, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}
func TestDoComputeDxyResult(t *testing.T) {
	val, err := DoComputeDxyResult(1.1, "1", "0", "0", "300", "", "", 2, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

func TestGetFloatValue(t *testing.T) {
	var p1 uint16 = 16481
	var p2 uint16 = 60293
	v := GetFloatValue(p1, p2)
	fmt.Println(v)
}

func TestGetFloatValue2(t *testing.T) {
	var p1 uint16 = 16481
	var p2 uint16 = 60293
	v := GetFloatValue2(p1, p2)
	fmt.Println(v)
}

func TestRing(t *testing.T) {
	ring := ring.New(3)

	ring.Value = 1
	ring.Move(1)
	ring.Value = 2
	fmt.Println(ring.Value)
	ring.Move(1)
	ring.Value = 3
	fmt.Println(ring.Value)
	ring.Move(1)
	ring = ring.Next()
	fmt.Println(ring.Value)

	ring.Do(func(p interface{}) {
		fmt.Println(p)
	})

}

func TestGetFloat32Value(t *testing.T) {
	formats := []string{BigEndian, BigEndianSwap, LittleEndian, LittleEndianSwap}
	src := []uint16{16481, 60293}
	for _, f := range formats {
		v := GetFloat32Value(src[0], src[1], f)
		fmt.Printf("Src=%d,%d Format=%s Value=%.3f\n", src[0], src[1], f, v)
	}
}

func TestGetFloat64Value(t *testing.T) {
	formats := []string{BigEndian, BigEndianSwap, LittleEndian, LittleEndianSwap}
	src := []uint16{16481, 60293, 16481, 60293}
	for _, f := range formats {
		v := GetFloat64Value(src[0], src[1], src[2], src[3], f)
		fmt.Printf("Src=%d,%d,%d,%d Format=%s Value=%.3f\n", src[0], src[1], src[2], src[3], f, v)
	}
}

func TestTcpVlidate(t *testing.T) {
	tv := NewTcpVlidate()
	fmt.Println(tv.Len(1))
	tv.SetLast(1, 1.2)
	tv.SetLast(1, 1.1)
	tv.SetLast(1, 1.2)
	fmt.Println(tv.Len(1))
	fmt.Println(tv.Sum(1))

}
