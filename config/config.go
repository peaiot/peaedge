package config

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/toughstruct/peaedge/common"
	"gopkg.in/yaml.v2"
)

type SysConfig struct {
	Appid      string `yaml:"appid"`
	Location   string `yaml:"location"`
	Workdir    string `yaml:"workdir"`
	DBFile     string `yaml:"dbfile"`
	SyslogAddr string `yaml:"syslog_addr"`
	Version    string `yaml:"version"`
	Debug      bool   `yaml:"debug"`
}

type WebConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	TlsPort int    `yaml:"tls_port"`
	Secret  string `yaml:"jwt_secret"`
	Debug   bool   `yaml:"debug"`
}

type ModbusConfig struct {
	TcpAddr  string `yaml:"tcp_addr"`
	RtuAddr  string `yaml:"rtu_addr"`
	BaudRate int    `yaml:"baud_rate"`
	DataBits int    `yaml:"data_bits"`
	StopBits int    `yaml:"stop_bits"`
	Parity   string `yaml:"parity"`
	Timeout  int    `yaml:"timeout"`
	Debug    bool   `yaml:"debug"`
}

type DataConfig struct {
	RtdSave string `yaml:"rtd_save"`
}

type AppConfig struct {
	System SysConfig    `yaml:"system"`
	Web    WebConfig    `yaml:"web"`
	Modbus ModbusConfig `yaml:"modbus"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}

func (c *AppConfig) GetPublicDir() string {
	return path.Join(c.System.Workdir, "public")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.System.Workdir, "private")
}

func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "public"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "private"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "backup"), 0644)
}

func setEnvValue(name string, val *string) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue
	}
}

func setEnvBoolValue(name string, val *bool) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue == "true" || evalue == "1" || evalue == "on"
	}
}

func setEnvInt64Value(name string, val *int64) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = p
	}
}
func setEnvIntValue(name string, val *int) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = int(p)
	}
}

var DefaultBssConfig = &AppConfig{
	System: SysConfig{
		Appid:      "peaedge1",
		Location:   "Asia/Shanghai",
		Workdir:    "/var/peaedge",
		DBFile:     "/var/peaedge/peaedge.db",
		SyslogAddr: "",
		Version:    "latest",
		Debug:      true,
	},
	Web: WebConfig{
		Host:   "0.0.0.0",
		Port:   1850,
		Secret: "9b6de5cc-0731-edge-peax-0f568ac9da37",
	},
	Modbus: ModbusConfig{
		TcpAddr: "0.0.0.0:1502",
		RtuAddr: "/dev/null",
		Debug:   true,
	},
}

func LoadConfig(cfile string) *AppConfig {
	// 开发环境首先查找当前目录是否存在自定义配置文件
	if cfile == "" {
		cfile = "peaedge.yml"
	}
	if !common.FileExists(cfile) {
		cfile = "/etc/peaedge.yml"
	}
	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(ioutil.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultBssConfig
	}
	setEnvValue("PEAEDGE_SYSTEM_WORKER_DIR", &cfg.System.Workdir)
	setEnvValue("PEAEDGE_SYSTEM_DBFILE", &cfg.System.DBFile)

	setEnvBoolValue("PEAEDGE_SYSTEM_DEBUG", &cfg.System.Debug)
	setEnvValue("PEAEDGE_SYSLOG_HOST", &cfg.System.SyslogAddr)
	setEnvValue("PEAEDGE_WEB_HOST", &cfg.Web.Host)
	setEnvValue("PEAEDGE_WEB_SECRET", &cfg.Web.Secret)
	setEnvIntValue("PEAEDGE_WEB_PORT", &cfg.Web.Port)
	setEnvBoolValue("PEAEDGE_WEB_DEBUG", &cfg.Web.Debug)

	return cfg
}
