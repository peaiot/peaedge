package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	_ "time/tzdata"

	"github.com/toughstruct/peaedge/apiserver"
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/common/installer"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/jobs"
	"github.com/toughstruct/peaedge/mqttc"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// Command line definition
var (
	h         = flag.Bool("h", false, "help usage")
	x         = flag.Bool("x", false, "debug pprof mode")
	showVer   = flag.Bool("v", false, "show version")
	conffile  = flag.String("c", "", "config yml file")
	initdb    = flag.Bool("initdb", false, "run initdb")
	install   = flag.Bool("install", false, "run install")
	uninstall = flag.Bool("uninstall", false, "run uninstall")
	initcfg   = flag.Bool("initcfg", false, "write default config > /etc/peaedge.yml")
)

// PrintVersion Print version information
func PrintVersion() {
	fmt.Fprintln(os.Stdout, assets.BuildInfo)
}

func printHelp() {
	if *h {
		ustr := fmt.Sprintf("bss Usage:bss -h\nOptions:")
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *showVer {
		PrintVersion()
		os.Exit(0)
	}

	if *x {
		go func() {
			log.Info(http.ListenAndServe("localhost:16060", nil))
		}()
	}

	_config := config.LoadConfig(*conffile)
	_config.InitDirs()

	printHelp()

	// 初始化配置文件
	if *initcfg {
		err := installer.InitConfig(_config)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// 安装为系统服务
	if *install {
		err := installer.Install(_config)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// 卸载
	if *uninstall {
		installer.Uninstall()
		return
	}

	// 完整初始化数据库
	if *initdb {
		app.Init(_config)
		app.Initdb()
		return
	}

	// 根据依赖关系按顺序进行初始化
	// 1-应用全局初始化
	app.Init(_config)
	app.Migrate(_config.System.Debug)
	// 2-任务调度初始化
	jobs.Init()
	defer app.OnExit()

	g.Go(func() error {
		log.Info("Start Web server ...")
		return apiserver.Listen()
	})

	g.Go(func() error {
		log.Info("Start matt client ...")
		return mqttc.StartDaemon()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
