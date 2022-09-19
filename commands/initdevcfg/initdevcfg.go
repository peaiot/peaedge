package main

import (
	"os"

	"github.com/toughstruct/peaedge/config"
	"gopkg.in/yaml.v2"
)

// 初始化一个本地开发配置文件

func main() {
	bs, err := yaml.Marshal(config.DefaultBssConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile("peaedge.yml", bs, 777)
}
