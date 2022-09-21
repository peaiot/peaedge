package app

import (
	slog "log"
	"os"
	"runtime/debug"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/common/log"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	Config *config.AppConfig
	DB     *gorm.DB
	// Cache  *freecache.Cache
)

// Init Global initialization call
func Init(cfg *config.AppConfig) {
	Config = cfg
	setupTimeZone()
	setupLogger()
	// Cache = freecache.NewCache(32 * 1024 * 1024)
	var err error
	DB, err = getGormDB()
	common.Must(err)
	log.Infof("app init done")
}
func setupTimeZone() {
	tz := Config.System.Location
	if tz == "" {
		tz = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Error(err)
		return
	}
	time.Local = loc
}

// Initialization log
func setupLogger() {
	level := logging.INFO
	if Config.System.Debug {
		level = logging.DEBUG
	}
	log.SetupLog(level, Config.System.SyslogAddr, Config.GetLogDir(), Config.System.Appid)
}

func getGormDB() (*gorm.DB, error) {
	pool, err := gorm.Open(sqlite.Open(Config.System.DBFile), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.New(
			slog.New(log.Stdlog{}, "\r\n", slog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				// LogLevel:                  common.If(Config.System.Debug, logger.Info, logger.Silent).(logger.LogLevel), // Log level
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,          // Disable color
			},
		),
	})
	if err != nil {
		return nil, err
	}
	sdb, _ := pool.DB()
	sdb.SetMaxOpenConns(1)
	return pool, nil
}

func OnExit() {
	if DB != nil {
		sdb, _ := DB.DB()
		_ = sdb.Close()
	}
}

func Migrate(track bool) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			if os.Getenv("GO_DEGUB_TRACE") != "" {
				debug.PrintStack()
			}
			err2, ok := err1.(error)
			if ok {
				err = err2
			}
		}
	}()
	if track {
		_ = DB.Debug().Migrator().AutoMigrate(models.Tables...)
	} else {
		_ = DB.Migrator().AutoMigrate(models.Tables...)
	}
	return nil
}

func Initdb() {
	_ = DB.Migrator().DropTable(models.Tables...)
	_ = DB.Migrator().AutoMigrate(models.Tables...)
}

func Debug() bool {
	return Config.System.Debug
}
