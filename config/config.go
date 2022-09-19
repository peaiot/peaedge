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
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	TlsPort   int    `yaml:"tls_port"`
	JwtSecret string `yaml:"jwt_secret"`
	Debug     bool   `yaml:"debug"`
}

type MqttConfig struct {
	Broker   string `yaml:"broker"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Debug    bool   `yaml:"debug"`
}

type AppConfig struct {
	System SysConfig  `yaml:"system"`
	Mqtt   MqttConfig `yaml:"mqtt"`
	Web    WebConfig  `yaml:"web"`
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

func (c *AppConfig) GetResourceDir() string {
	return path.Join(c.System.Workdir, "resource")
}

func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "public"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "private"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "resource"), 0755)
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
		Appid:      "peaedge-1",
		Location:   "Asia/Shanghai",
		Workdir:    "/var/peaedge",
		DBFile:     "/var/peaedge/peaedge.db",
		SyslogAddr: "",
		Version:    "latest",
		Debug:      true,
	},
	Mqtt: MqttConfig{
		Broker: "tcp://",
	},
	Web: WebConfig{
		Host:      "0.0.0.0",
		Port:      1850,
		JwtSecret: "9b6de5cc-0731-edge-peax-0f568ac9da37",
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

	setEnvValue("PEAEDGE_MQTT_BROKER", &cfg.Mqtt.Broker)
	setEnvValue("PEAEDGE_MQTT_USERNAME", &cfg.Mqtt.Username)
	setEnvValue("PEAEDGE_MQTT_PASSWORD", &cfg.Mqtt.Password)
	setEnvBoolValue("PEAEDGE_MQTT_DEBUG", &cfg.Mqtt.Debug)

	setEnvBoolValue("PEAEDGE_SYSTEM_DEBUG", &cfg.System.Debug)
	setEnvValue("PEAEDGE_SYSLOG_HOST", &cfg.System.SyslogAddr)
	setEnvValue("PEAEDGE_WEB_HOST", &cfg.Web.Host)
	setEnvValue("PEAEDGE_WEB_SECRET", &cfg.Web.JwtSecret)
	setEnvIntValue("PEAEDGE_WEB_PORT", &cfg.Web.Port)
	setEnvBoolValue("PEAEDGE_WEB_DEBUG", &cfg.Web.Debug)

	return cfg
}