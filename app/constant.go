package app

const (
	DataFlagSuccess     = "N" // 在线监控(监测)仪器仪表工作正常
	DataFlagStop        = "F" // 在线监控(监测)仪器仪表停运
	DataFlagMaintenance = "M" // 在线监控(监测)仪器仪表处于维护期间产生的数据
	DataFlagManuals     = "S" // 手工输入的设定值
	DataFlagFailure     = "D" //  在线监控(监测)仪器仪表故障
	DataFlagAdjust      = "C" //  在线监控(监测)仪器仪表处于校准状态
	DataFlagOver        = "T" //  在线监控(监测)仪器仪表采样数值超过测量上限
	DataFlagNetBreak    = "B" //  在线监控(监测)仪器仪表与数采仪通讯异常
)
