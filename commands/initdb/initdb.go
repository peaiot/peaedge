package main

import (
	"github.com/toughstruct/peaedge/app"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/models"
)

func main() {
	app.Init(config.LoadConfig(""))
	_ = app.GormDB.Migrator().AutoMigrate(models.Tables...)
}
