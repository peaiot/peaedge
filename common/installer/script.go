package installer

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/config"

	"gopkg.in/yaml.v2"
)

var InstallScript = `#!/bin/bash -x
mkdir -p /var/peaedge
chmod -R 777 /var/peaedge
install -m 777 ./peaedge /usr/local/bin/peaedge
test -d /usr/lib/systemd/system || mkdir -p /usr/lib/systemd/system
test -d /etc/sysconfig || mkdir -p /etc/sysconfig
test -f /etc/sysconfig/peaedge || touch /etc/sysconfig/peaedge
cat>/usr/lib/systemd/system/peaedge.service<<EOF
[Unit]
Description=peaedge
After=network.target
StartLimitIntervalSec=0

[Service]
Restart=always
RestartSec=1
EnvironmentFile=/etc/sysconfig/peaedge
Environment=GODEBUG=x509ignoreCN=0
LimitNOFILE=65535
LimitNPROC=65535
User=root
ExecStart=/usr/local/bin/peaedge

[Install]
WantedBy=multi-user.target
EOF

chmod 600 /usr/lib/systemd/system/peaedge.service
systemctl enable peaedge && systemctl daemon-reload

`

func InitConfig(config *config.AppConfig) error {
	// config.NBI.JwtSecret = common.UUID()
	cfgstr, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("/etc/peaedge.yaml", cfgstr, 0644)
}

func Install(config *config.AppConfig) error {
	if !common.FileExists("/etc/peaedge.yaml") {
		_ = InitConfig(config)
	}
	script := strings.ReplaceAll(InstallScript, "/var/peaedge", config.System.Workdir)
	cmd := "/usr/local/bin/peaedge"
	script = strings.ReplaceAll(InstallScript, "/usr/local/bin/peaedge", cmd)
	_ = ioutil.WriteFile("/tmp/peaedge_install.sh", []byte(script), 0777)

	// 创建用户&组
	if err := exec.Command("/bin/bash", "/tmp/peaedge_install.sh").Run(); err != nil {
		return err
	}
	return os.Remove("/tmp/peaedge_install.sh")
}

func Uninstall() {
	_ = os.Remove("/etc/peaedge.yaml")
	_ = os.Remove("/usr/lib/systemd/system/peaedge.service")
	_ = os.Remove("/usr/local/bin/peaedge")
}
