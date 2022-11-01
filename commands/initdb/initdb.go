package main

import (
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/config"
)

func main() {
	app.Init(config.LoadConfig(""))
	app.Initdb()
}
