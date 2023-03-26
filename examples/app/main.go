package main

import (
	"github.com/netr/napi"
	"github.com/netr/napi/examples/app/web/ctrl"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
)

var (
	db *gorm.DB
)

func main() {
	var err error
	s := napi.NewServer(
		napi.DefaultFiberConfig("test_app"),
		napi.WithCatchAll(),
	).
		Port(1338).UseBaseMiddlewares().
		UsePrometheus().UsePprof().UseHealth().
		UseDefaultCORS().UseDefaultLogger().UseDefaultLimiter()

	db, err = newGormDB()
	handleErr(err)

	ctrl.NewRoutes(s.App()).Setup(db)

	s.Run()
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func newGormDB() (*gorm.DB, error) {
	gormDb, err := gorm.Open(sqlite.Open("sqlite.db"),
		&gorm.Config{
			Logger: glog.Default.LogMode(glog.Silent),
		},
	)
	if err != nil {
		return nil, err
	}
	migrateDB(gormDb)

	return gormDb, nil
}

func migrateDB(db *gorm.DB, migrations ...interface{}) {
	_ = db.AutoMigrate(migrations...)
}
