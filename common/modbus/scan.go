package modbus

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ScanPorts() []string {

	switch runtime.GOOS {
	case "darwin":
		bs, _ := exec.Command("/bin/sh", "-c", "ls /dev/{tty,cu}.*").Output()
		ss := strings.Split(string(bs), "\n")
		for i, v := range ss {
			ss[i] = strings.TrimSpace(v)
		}
		return ss
	case "linux":
		ports, _ := ioutil.ReadDir("/dev/serial/by-id/")
		numOfPorts := len(ports)
		if numOfPorts == 0 {
			return nil
		}
		_ports := make([]string, 0)
		for _, port := range ports {
			portPath, _ := filepath.EvalSymlinks("/dev/serial/by-id/" + port.Name())
			_ports = append(_ports, strings.TrimSpace(portPath))
		}
		return _ports
	case "windows":
		out, _ := exec.Command("powershell", "Get-WmiObject Win32_SerialPort", "|", "Select DeviceID -ExpandProperty DeviceID").Output()
		ss := strings.Split(string(out), "\n")
		for i, v := range ss {
			ss[i] = strings.TrimSpace(v)
		}
		return ss
	default:
		return nil
	}
	return nil
}
