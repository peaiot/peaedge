package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peaiot/logging"
)

const ModuleSystem = "System"
const ModuleHj212 = "HjZIZ"
const ModuleMqttc = "Mqttc"
const ModuleModbus = "Modbus"

var Default = logging.MustGetLogger(ModuleSystem)
var Hj212 = logging.MustGetLogger(ModuleHj212)
var Mqttc = logging.MustGetLogger(ModuleMqttc)
var Modbus = logging.MustGetLogger(ModuleModbus)

func SetupLog(level logging.Level, syslogaddr string, logdir string) {
	var format = logging.MustStringFormatter(
		`%{color} %{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	Backends := make([]logging.Backend, 0)
	Backends = append(Backends, logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))
	bs := SetupSyslog(level, syslogaddr, ModuleSystem)
	bf := FileSyslog(level, logdir, ModuleSystem)
	if bs != nil {
		Backends = append(Backends, bs)
	}
	if bf != nil {
		Backends = append(Backends, bf)
	}
	logging.SetBackend(Backends...)
	logging.SetLevel(level, ModuleSystem)

	hjbf := FileSyslog(level, logdir, ModuleHj212)
	if hjbf != nil {
		Hj212.SetBackend(hjbf)
	}
	mqbf := FileSyslog(level, logdir, ModuleMqttc)
	if mqbf != nil {
		Mqttc.SetBackend(mqbf)
	}
	mbbf := FileSyslog(level, logdir, ModuleModbus)
	if mbbf != nil {
		Modbus.SetBackend(mbbf)
	}

}

func clearLogs(logsdir string, prefix string) {
	daydirs, err := ioutil.ReadDir(logsdir)
	if err != nil {
		Default.Errorf("read day logs dir error, %s", err.Error())
		return
	}

	for _, item := range daydirs {
		if !item.IsDir() && strings.HasPrefix(item.Name(), prefix) && item.ModTime().Before(time.Now().Add(-(time.Hour * 24 * 7))) {
			fpath := filepath.Join(logsdir, item.Name())
			err = os.Remove(fpath)
			if err != nil {
				Default.Errorf("remove logfile %s error", fpath)
			}
		}
	}
}

func FileSyslog(level logging.Level, logdir string, module string) logging.LeveledBackend {
	if logdir == "N/A" {
		return nil
	}
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
	)

	logfile, err := NewFile(filepath.Join(logdir, module+"-daily-2006-01-02.log"), func(path string, didRotate bool) {
		fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
		if !didRotate {
			return
		}
		// process just closed file e.g. upload to backblaze storage for backup
		go clearLogs(logdir, module+"-daily-")
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backendFile := logging.NewLogBackend(logfile, "", 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backend2Formatter := logging.NewBackendFormatter(backendFile, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(level, module)
	return backend1Leveled
}

var (
	Error  = Default.Error
	Errorf = Default.Errorf
	Info   = Default.Info
	Infof  = Default.Infof
	Fatal  = Default.Fatal
	Fatalf = Default.Fatalf
	Debug  = Default.Debug
	Debugf = Default.Debugf

	IsDebug = func() bool {
		return Default.IsEnabledFor(logging.DEBUG)
	}
)
