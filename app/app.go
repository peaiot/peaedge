package app

import (
	slog "log"
	"os"
	"path"
	"runtime/debug"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nakabonne/tstorage"
	"github.com/peaiot/logging"
	"github.com/toughstruct/peaedge/common"
	"github.com/toughstruct/peaedge/config"
	"github.com/toughstruct/peaedge/log"
	"github.com/toughstruct/peaedge/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	appConfig *config.AppConfig
	gormDB    *gorm.DB
	tsdb      tstorage.Storage
	isInit    bool
	// cache  *freecache.Cache
)

// Init Global initialization call
func Init(cfg *config.AppConfig) {
	appConfig = cfg
	setupTimeZone()
	setupLogger()
	setupTsStorage()
	// Cache = freecache.NewCache(32 * 1024 * 1024)
	var err error
	gormDB, err = getGormDB()
	common.Must(err)
	isInit = true
	log.Infof("app init done")
}

func setupTsStorage() {
	tsdb, _ = tstorage.NewStorage(
		tstorage.WithPartitionDuration(time.Hour),
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithRetention(24*time.Hour),
		tstorage.WithDataPath(path.Join(appConfig.System.Workdir, "tstorage")),
		tstorage.WithWALBufferedSize(0),
		tstorage.WithWriteTimeout(60*time.Second),
	)
}

func setupTimeZone() {
	tz := appConfig.System.Location
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
	if appConfig.System.Debug {
		level = logging.DEBUG
	}
	log.SetupLog(level, appConfig.System.SyslogAddr, appConfig.GetLogDir())
}

func getGormDB() (*gorm.DB, error) {
	pool, err := gorm.Open(sqlite.Open(appConfig.System.DBFile), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.New(
			slog.New(log.Default, "\r\n", slog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				// LogLevel:                  common.If(Config.System.IsDebug, logger.Info, logger.Silent).(logger.LogLevel), // Log level
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
	if gormDB != nil {
		sdb, _ := gormDB.DB()
		_ = sdb.Close()
	}

	if tsdb != nil {
		_ = tsdb.Close()
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
		_ = gormDB.Debug().Migrator().AutoMigrate(models.Tables...)
	} else {
		_ = gormDB.Migrator().AutoMigrate(models.Tables...)
	}
	return nil
}

func Initdb() {
	_ = gormDB.Migrator().DropTable(models.Tables...)
	_ = gormDB.Migrator().AutoMigrate(models.Tables...)
}

func Config() *config.AppConfig {
	return appConfig
}
func DB() *gorm.DB {
	return gormDB
}

func TsDB() tstorage.Storage {
	return tsdb
}

func IsDebug() bool {
	return appConfig.System.Debug
}

func IsInit() bool {
	return isInit
}
