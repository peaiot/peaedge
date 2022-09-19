package e212

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	HEAD = "##"
	TAIL = "\r\n"

	// 提取现场机时间 device <-- server
	RequestCNDeviceTime = 1011
	// 上传现场机时间 device --> server
	UploadCNDeviceTime = 1011
	// 设置现场机时间 device <-- server
	RequestCNSetDeviceTime = 1012
	// 现场机时间校准请求 device --> server
	NotifyCNDeviceTimeAdjust = 1013
)

var (
	DataRegexp    = regexp.MustCompile(`(QN=\w+);(ST=\w+);(CN=\w+);(PW=\w*);(MN=\w+);(Flag=\w+);(PNUM=\w+)?;?(PNO=\w+)?;?CP=&&(.+)&&`)
	PNUMPNORegexp = regexp.MustCompile(`(PNUM=\w+);(PNO=\w+);`)
)

// 封装key=value
func warpKeyValue(key string, value string) string {
	var buf strings.Builder
	buf.WriteString(key)
	buf.WriteString("=")
	buf.WriteString(value)
	return buf.String()
}

// 数据段定义
type DataSegment struct {
	_QN   string
	_ST   string
	_CN   string
	_PW   string
	_MN   string
	_Flag string
	_PNUM string
	_PNO  string
	_CP   string
}

// 从字段属性值构建数据段报文
func NewDataSegment(QN string, ST string, CN string, PW string, MN string, flag string, CP string) *DataSegment {
	ds := &DataSegment{}
	ds.SetQNValue(QN)
	ds.SetSTValue(ST)
	ds.SetCNValue(CN)
	ds.SetPWValue(PW)
	ds.SetMNValue(MN)
	ds.SetFlagValue(flag)
	ds.SetCpValue(CP)
	return ds
}

// 解析数据段原始报文
func ParseDataSegmentFrom(data string) (*DataSegment, error) {
	// var pos = 0
	ds := &DataSegment{}
	attrs := DataRegexp.FindStringSubmatch(data)
	if attrs == nil {
		return nil, errors.New("解析数据段错误")
	}
	ds._QN = attrs[1]
	ds._ST = attrs[2]
	ds._CN = attrs[3]
	ds._PW = attrs[4]
	ds._MN = attrs[5]
	ds._Flag = attrs[6]
	ds._PNUM = attrs[7]
	ds._PNO = attrs[8]
	ds._CP = "CP=&&" + attrs[9] + "&&"
	pattrs := PNUMPNORegexp.FindStringSubmatch(data)
	if pattrs != nil && len(pattrs) == 3 {
		ds._PNUM = pattrs[1]
		ds._PNO = pattrs[2]
	}
	return ds, nil
}

// 请求编码 QN
func (s *DataSegment) SetQNValue(val string) {
	s._QN = warpKeyValue("QN", val)
}

func (s *DataSegment) GetQNValue() string {
	return s._QN[3:]
}

// 系统编码 ST
func (s *DataSegment) SetSTValue(val string) {
	s._ST = warpKeyValue("ST", val)
}

func (s *DataSegment) GetSTValue() string {
	return s._ST[3:]
}

// 命令编码 CN
func (s *DataSegment) SetCNValue(val string) {
	s._CN = warpKeyValue("CN", val)
}

func (s *DataSegment) GetCNValue() string {
	return s._CN[3:]
}

// 访问密码
func (s *DataSegment) SetPWValue(val string) {
	s._PW = warpKeyValue("PW", val)
}

func (s *DataSegment) GetPWValue() string {
	return s._PW[3:]
}

// 设备唯一标识 MN
func (s *DataSegment) SetMNValue(val string) {
	s._MN = warpKeyValue("MN", val)
}

func (s *DataSegment) GetMNValue() string {
	return s._MN[3:]
}

// 拆分包及应答标志
func (s *DataSegment) SetFlagValue(val string) {
	s._Flag = warpKeyValue("Flag", val)
}

func (s *DataSegment) GetFlagValue() string {
	return s._Flag[5:]
}

func (s *DataSegment) GetPNUMValue() string {
	if s._PNUM != "" {
		return s._PNUM[5:]
	}
	return ""
}

func (s *DataSegment) GetPNOValue() string {
	if s._PNO != "" {
		return s._PNO[4:]
	}
	return ""
}

// 设置 CP 数据区值， 参数不要&&
func (s *DataSegment) SetCpValue(val string) {
	var buf strings.Builder
	buf.WriteString("CP=&&")
	buf.WriteString(val)
	buf.WriteString("&&")
	s._CP = buf.String()
}

func (s *DataSegment) GetCPValue() string {
	return s._CP
}

// 将数据区解析为一个字典数据结构
func (s *DataSegment) GetCpAttrMap() map[string]map[string]interface{} {
	var _cpval = s._CP[5 : len(s._CP)-2]
	var result = map[string]map[string]interface{}{}
	result["attrs"] = make(map[string]interface{})
	var groups = strings.Split(_cpval, ";")
	for _, group := range groups {
		var attrs = strings.Split(group, ",")
		for _, attr := range attrs {
			nameAndvalue := strings.Split(attr, "=")
			if nameAndvalue != nil && len(nameAndvalue) == 2 {
				name := nameAndvalue[0]
				value := nameAndvalue[1]
				factorAndMetrics := strings.Split(name, "-")
				if factorAndMetrics != nil && len(factorAndMetrics) == 2 {
					factor := factorAndMetrics[0]
					metrics := factorAndMetrics[1]
					if _, ok := result[factor]; !ok {
						result[factor] = map[string]interface{}{}
					}
					result[factor][metrics] = value
				} else {
					result["attrs"][name] = value
				}
			}
		}
	}
	return result
}

// 返回文档格式， 便于存储
func (s *DataSegment) Document() map[string]interface{} {
	var result = map[string]interface{}{}
	result["QN"] = s.GetQNValue()
	result["ST"] = s.GetSTValue()
	result["CN"] = s.GetCNValue()
	result["PW"] = s.GetPWValue()
	result["MN"] = s.GetMNValue()
	result["Flag"] = s.GetFlagValue()
	result["PNUM"] = s.GetPNUMValue()
	result["PNO"] = s.GetPNOValue()
	result["CP"] = s.GetCpAttrMap()
	return result
}

// 数据区编码
func (s *DataSegment) Encode() string {
	var buf strings.Builder
	buf.WriteString(s._QN)
	buf.WriteString(";")
	buf.WriteString(s._ST)
	buf.WriteString(";")
	buf.WriteString(s._CN)
	buf.WriteString(";")
	buf.WriteString(s._PW)
	buf.WriteString(";")
	buf.WriteString(s._MN)
	buf.WriteString(";")
	buf.WriteString(s._Flag)
	buf.WriteString(";")
	if s._PNUM != "" {
		buf.WriteString(s._PNUM)
		buf.WriteString(";")
	}
	if s._PNUM != "" {
		buf.WriteString(s._PNO)
		buf.WriteString(";")
	}
	buf.WriteString(s._CP)
	return buf.String()
}

type E212Packet struct {
	Len     string
	Segment *DataSegment
	Crc     string
}

// 从原始报文字符串解析
func ParsePacketFrom(data string) (*E212Packet, error) {
	p := &E212Packet{}
	if !strings.HasPrefix(data, "##") {
		return nil, errors.New("")
	}
	p.Len = data[2:6]
	intlen, _ := strconv.Atoi(p.Len)
	segmentStr := data[6 : intlen+6]
	_Segment, err := ParseDataSegmentFrom(segmentStr)
	if err != nil {
		return nil, err
	}
	p.Segment = _Segment
	p.Crc = strings.TrimSpace(data[6+intlen:])
	// 校验CRC
	actcrc := Crc(segmentStr)
	if p.Crc != actcrc {
		return nil, errors.New(fmt.Sprintf("Crc check error, Expect %s actual %s", p.Crc, actcrc))
	}
	return p, nil
}

// 根据数据段构建报文
func NewE212Packet(segment *DataSegment) *E212Packet {
	return &E212Packet{Segment: segment}
}

// 编码报文
func (p *E212Packet) Encode() (string, error) {
	var segmentStr = p.Segment.Encode()
	var seglen = len(segmentStr)
	if seglen > 1024 {
		return "", errors.New("DataSegment length > 1024")
	}
	p.Len = fmt.Sprintf("%04d", seglen)

	var build strings.Builder
	build.WriteString(HEAD)
	build.WriteString(p.Len)
	build.WriteString(segmentStr)
	build.WriteString(Crc(segmentStr))
	build.WriteString(TAIL)
	return build.String(), nil
}
