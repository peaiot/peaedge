package e212

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

var cp = `DataTime=20200421212100;a01013-SampleTime=,a01013-Rtd=0.000,a01013-Flag=D;a01012-SampleTime=,a01012-Rtd=0.000,a01012-Flag=D;a01014-SampleTime=,a01014-Rtd=0.000,a01014-Flag=D;a01011-SampleTime=,a01011-Rtd=0.000,a01011-Flag=D;a24088-SampleTime=,a24088-Rtd=0.000,a24088-Flag=D;a05002-SampleTime=,a05002-Rtd=0.000,a05002-Flag=D`

func TestEncode(t *testing.T) {
	var seg = NewDataSegment("20200421212154904", "31", "2011", "123456", "MYHJ00000000000080011121", "4", cp)
	var p = NewE212Packet(seg)
	fmt.Println(p.Encode())
}

func TestCovBin(t *testing.T) {
	fmt.Println()
}

func TestParseFrom(t *testing.T) {
	data := `##0451QN=20200428120745080;ST=31;CN=2011;PW=123456;MN=MYHJ00000000000080011128;Flag=4;CP=&&DataTime=20200428120300;a24088-Rtd=0.000,a24088-Flag=D;a24087-Rtd=0.000,a24087-Flag=D;a25005-Rtd=0.000,a25005-Flag=D;a25003-Rtd=0.000,a25003-Flag=D;a25002-Rtd=0.000,a25002-Flag=D;a05002-Rtd=0.000,a05002-Flag=D;a01014-Rtd=0.000,a01014-Flag=D;a01013-Rtd=0.000,a01013-Flag=D;a01012-Rtd=0.000,a01012-Flag=D;a01011-Rtd=0.000,a01011-Flag=D;a00000-Rtd=0.000,a00000-Flag=D&&0000`
	dlen := data[2:6]
	fmt.Println("datalen:", dlen)
	intlen, _ := strconv.Atoi(dlen)
	segmentStr := data[6 : intlen+6]
	fmt.Println("segmentStr:", segmentStr)
	packet, err := ParsePacketFrom(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%+v", packet.Segment.Document()))
	v, _ := json.MarshalIndent(packet.Segment.Document(), "", "\t")
	fmt.Println(string(v))

}

func TestParseValidFrom(t *testing.T) {
	data := `##0214QN=20200721093700136;ST=22;CN=2011;PW=123456;MN=MYHJ00000000000080011129;Flag=4;CP=&&DataTime=20200721093700;a99054-Rtd=190.000;a01001-Rtd=0.000;a01002-Rtd=0.000;a01006-Rtd=0.000;a01007-Rtd=0.000;a01008-Rtd=0.000&&33EC`
	dlen := data[2:6]
	fmt.Println("datalen:", dlen)
	fmt.Println("srclen:", len(data))
	intlen, _ := strconv.Atoi(dlen)
	segmentStr := data[6 : intlen+6]
	fmt.Println("segmentStr:", segmentStr)
	packet, err := ParsePacketFrom(data)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(fmt.Sprintf("%+v", packet.Segment.Document()))
	v, _ := json.MarshalIndent(packet.Segment.Document(), "", "\t")
	fmt.Println(string(v))

}

func TestReg(t *testing.T) {
	flysnowRegexp := regexp.MustCompile(`(Flag=\d+);`)
	params := flysnowRegexp.FindStringSubmatch("Flag=4;CP")

	for _, param := range params {
		fmt.Println(param)
	}
}

func TestReg2(t *testing.T) {
	_PNUMPNORegexp := regexp.MustCompile(`(PNUM=\w+);(PNO=\w+);`)
	params := _PNUMPNORegexp.FindStringSubmatch("PNU0M=4;P0NO=2;")

	for _, param := range params {
		fmt.Println(param)
	}
}

func TestReg3(t *testing.T) {
	ss := `##0131QN=20200428011552615;ST=31;CN=2081;PW=;MN=MYHJ00000000000080011128;Flag=4;CP=&&DataTime=20200428011552;RestartTime=20200425120954&&47C0`
	ss2 := `##0131QN=20200428102058322;ST=31;CN=2081;PW=123456;MN=MYHJ00000000000080011128;Flag=4;CP=&&DataTime=20200428102058;RestartTime=20200428100633&&47c0`
	attrs := DataRegexp.FindStringSubmatch(ss)
	for _, param := range attrs {
		fmt.Println(param)
	}
	fmt.Println("--------------")
	attrs2 := DataRegexp.FindStringSubmatch(ss2)
	for _, param := range attrs2 {
		fmt.Println(param)
	}
}

func TestMsg(t *testing.T) {
	data := `##0931QN=20201029122005355;ST=31;CN=2051;PW=123456;MN=WJTHGY91320509703697203A;Flag=7;PNUM=2;PNO=1;CP=&&DataTime=20201029121000;a00000-Min=0.000,a00000-Max=0.000,a00000-Avg=0.000,a00000-Cou=0.000,a00000-Flag=B;a01011-Min=0.000,a01011-Max=0.000,a01011-Avg=0.000,a01011-Cou=0.000,a01011-Flag=B;a01012-Min=0.000,a01012-Max=0.000,a01012-Avg=0.000,a01012-Cou=0.000,a01012-Flag=B;a01013-Min=0.000,a01013-Max=0.000,a01013-Avg=0.000,a01013-Cou=0.000,a01013-Flag=B;a01014-Min=0.000,a01014-Max=0.000,a01014-Avg=0.000,a01014-Cou=0.000,a01014-Flag=B;a19001-Min=0.000,a19001-Max=0.000,a19001-Avg=0.000,a19001-Cou=0.000,a19001-Flag=B;a21002-Min=0.000,a21002-Max=0.000,a21002-Avg=0.000,a21002-Cou=0.000,a21002-ZsMin=0.000,a21002-ZsMax=0.000,a21002-ZsAvg=0.000,a21002-ZsCou=0.000,a21002-Flag=B;a21005-Min=0.000,a21005-Max=0.000,a21005-Avg=0.000,a21005-Cou=0.000,a21005-ZsMin=0.000,a21005-ZsMax=0.000,a21005-ZsAvg=0.000,a21005-ZsCou=0.000,a21005-Flag=B&&AE41`
	fmt.Println(len(data))
	dlen := data[2:6]
	fmt.Println("datalen:", dlen)
	intlen, _ := strconv.Atoi(dlen)
	segmentStr := data[6 : intlen+6]
	fmt.Println("segmentStr:", segmentStr)
	packet, err := ParsePacketFrom(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%+v", packet.Segment.Document()))
	v, _ := json.MarshalIndent(packet.Segment.Document(), "", "\t")
	fmt.Println(string(v))
}
func TestMsgx(t *testing.T) {
	data := `##0494QN=20210521110212000;ST=31;CN=2011;PW=123456;MN=411420sibo0606;Flag=4;CP=&&DataTime=20210521110202;a01012-Rtd=32.42,a01012-ZsRtd=32.42,a01012-Flag=N;a05002-Rtd=1.45,a05002-ZsRtd=1.45,a05002-Flag=N;a01011-Rtd=0.54,a01011-ZsRtd=0.54,a01011-Flag=N;a01013-Rtd=103.79,a01013-ZsRtd=103.79,a01013-Flag=N;a24087-Rtd=39.00,a24087-ZsRtd=39.00,a24087-Flag=N;a01014-Rtd=2.53,a01014-ZsRtd=2.53,a01014-Flag=N;a24088-Rtd=37.55,a24088-ZsRtd=37.55,a24088-Flag=N;a00000-Rtd=3434,a00000-ZsRtd=3434,a00000-Flag=N&&8701`
	fmt.Println(len(data))
	dlen := data[2:6]
	fmt.Println("datalen:", dlen)
	intlen, _ := strconv.Atoi(dlen)
	segmentStr := data[6 : intlen+6]
	fmt.Println("segmentStr:", segmentStr)
	packet, err := ParsePacketFrom(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%+v", packet))
}

func TestMsgxx(t *testing.T) {
	data := "##0480QN=20210521110006000;ST=31;CN=2011;PW=123456;MN=411420sibo0606;Flag=4;CP=&&DataTime=20210521110000;a01014-Rtd=2.47,a01014-ZsRtd=2.47,a01014-Flag=N;a24088-Rtd=37.16,a24088-ZsRtd=37.16,a24088-Flag=N;a00000-Rtd=0,a00000-ZsRtd=0,a00000-Flag=N;a01011-Rtd=0,a01011-ZsRtd=0,a01011-Flag=N;a24087-Rtd=38.52,a24087-ZsRtd=38.52,a24087-Flag=N;a01012-Rtd=32.69,a01012-ZsRtd=32.69,a01012-Flag=N;a05002-Rtd=1.36,a05002-ZsRtd=1.36,a05002-Flag=N;a01013-Rtd=61.05,a01013-ZsRtd=61.05,a01013-Flag=N&&1040"
	fmt.Println(len(data))
	dlen := data[2:6]
	fmt.Println("datalen:", dlen)
	intlen, _ := strconv.Atoi(dlen)
	segmentStr := data[6 : intlen+6]
	fmt.Println("segmentStr:", segmentStr)
	packet, err := ParsePacketFrom(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%+v", packet))
}

func TestCrc(t *testing.T) {
	fmt.Println(fmt.Sprintf("%04s","0"))
}

func TestCrcgen(t *testing.T) {
	fmt.Println(Crc("QN=20200721093700136;ST=22;CN=2011;PW=123456;MN=MYHJ00000000000080011129;Flag=4;CP=&&DataTime=20200721093700;a99054-Rtd=190.000;a01001-Rtd=0.000;a01002-Rtd=0.000;a01006-Rtd=0.000;a01007-Rtd=0.000;a01008-Rtd=0.000&&"))
}


//segmentStr: QN=20210521110212000;ST=31;CN=2011;PW=123456;MN=411420sibo0606;Flag=4;CP=&&DataTime=20210521110202;a01012-Rtd=32.42,a01012-ZsRtd=32.42,a01012-Flag=N;a05002-Rtd=1.45,a05002-ZsRtd=1.45,a05002-Flag=N;a01011-Rtd=0.54,a01011-ZsRtd=0.54,a01011-Flag=N;a01013-Rtd=103.79,a01013-ZsRtd=103.79,a01013-Flag=N;a24087-Rtd=39.00,a24087-ZsRtd=39.00,a24087-Flag=N;a01014-Rtd=2.53,a01014-ZsRtd=2.53,a01014-Flag=N;a24088-Rtd=37.55,a24088-ZsRtd=37.55,a24088-Flag=N;a00000-Rtd=3434,a00000-ZsRtd=3434,a00000-Flag=N&&
//&{Len:0494 Segment:0xc0000f4000 Crc:8701}
//
//segmentStr: QN=20210521110006000;ST=31;CN=2011;PW=123456;MN=411420sibo0606;Flag=4;CP=&&DataTime=20210521110000;a01014-Rtd=2.47,a01014-ZsRtd=2.47,a01014-Flag=N;a24088-Rtd=37.16,a24088-ZsRtd=37.16,a24088-Flag=N;a00000-Rtd=0,a00000-ZsRtd=0,a00000-Flag=N;a01011-Rtd=0,a01011-ZsRtd=0,a01011-Flag=N;a24087-Rtd=38.52,a24087-ZsRtd=38.52,a24087-Flag=N;a01012-Rtd=32.69,a01012-ZsRtd=32.69,a01012-Flag=N;a05002-Rtd=1.36,a05002-ZsRtd=1.36,a05002-Flag=N;a01013-Rtd=61.05,a01013-ZsRtd=61.05,a01013-Flag=N&&
//&{Len:0480 Segment:0xc00012e000 Crc:1040}
