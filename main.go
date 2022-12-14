package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"runtime"
	_ "time/tzdata"

	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/assets"
	"github.com/toughstruct/peaedge/channels/httpc"
	"github.com/toughstruct/peaedge/channels/mqttc"
	"github.com/toughstruct/peaedge/channels/tcpc"
	"github.com/toughstruct/peaedge/common/installer"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/jobs"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/mbslave"
	"github.com/toughstruct/peaedge/webserver"
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
	initdb    = flag.Bool("initdb", false, "initdb")
	inittest  = flag.Bool("init-test", false, "init test data")
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

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	// go open("http://127.0.0.1:1850")
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

	// ?????????????????????
	if *initcfg {
		err := installer.InitConfig(_config)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// ?????????????????????
	if *install {
		err := installer.Install(_config)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// ??????
	if *uninstall {
		installer.Uninstall()
		return
	}

	// ????????????????????????
	if *initdb {
		app.Init(_config)
		app.Initdb()
		return
	}

	if *inittest {
		app.Init(_config)
		app.InitTestData()
		return
	}

	// ??????????????????????????????????????????
	// 1-?????????????????????
	app.Init(_config)
	_ = app.Migrate(*x)
	// 2-?????????????????????
	jobs.Init()
	defer app.OnExit()

	if err := mqttc.StartAll(); err != nil {
		log.Error(err)
	}

	if err := httpc.StartAll(); err != nil {
		log.Error(err)
	}

	if err := tcpc.StartAll(); err != nil {
		log.Error(err)
	}

	g.Go(func() error {
		log.Info("Start Web server ...")
		return webserver.Listen()
	})

	g.Go(func() error {
		log.Info("Start Modbus server ...")
		return mbslave.Listen()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
