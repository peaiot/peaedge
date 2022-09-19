package sysid

import (
	"net"

	"github.com/denisbrodbeck/machineid"
	"github.com/toughstruct/peaedge/common"
)

func GetSystemSid() string {
	sid, err := machineid.ProtectedID("peaedge")
	if err != nil {
		return GetFirstMacAddrSid()
	}
	return common.Md5Hash(sid)
}

func GetFirstMacAddrSid() string {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.HardwareAddr.String() != "" {
			return common.Md5Hash(inter.HardwareAddr.String())
		}
	}
	return common.Md5Hash("peaedge")
}
