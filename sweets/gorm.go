package sweets

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
	"reflect"
	"testing"
)

type GormSuite struct {
	db         *gorm.DB
	ranOnce    bool
	migrations []interface{}
}

// NewGormSuite is used to instantiate a new gorm.DB test suite. Typically called in SetupSuite().
// We can leverage log.Fatal here, since this method is only used in testing, removing verbosity from our test files.
func (suite *GormSuite) NewGormSuite(migrations ...interface{}) {
	db, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{
			Logger: glog.Default.LogMode(glog.Silent),
		})
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	err = db.AutoMigrate(migrations...)
	if err != nil {
		log.Fatal(err)
	}

	suite.migrations = migrations
	suite.ranOnce = false
	return
}

// RefreshDB will drop all your current migrations and re-migrate with a fresh database. Used in SetupTest().
// We can leverage log.Fatal here, since this method is only used in testing, removing verbosity from our test files.
func (suite *GormSuite) RefreshDB() {
	if suite.ranOnce {
		err := suite.db.Migrator().DropTable(suite.migrations...)
		if err != nil {
			log.Fatal(err)
		}

		err = suite.db.AutoMigrate(suite.migrations...)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		suite.ranOnce = true
	}

	return
}

// DB is a helper function to get the underlying *gorm.DB
func (suite *GormSuite) DB() *gorm.DB {
	return suite.db
}

// Sql is a helper function to get the underlying *sql.DB
func (suite *GormSuite) Sql() *sql.DB {
	if s, err := suite.db.DB(); err != nil {
		return nil
	} else {
		return s
	}
}

// ShutdownDB is a helper function to shut down the underlying *gorm.DB
func (suite *GormSuite) ShutdownDB() {
	if d, err := suite.DB().DB(); err != nil {
		return
	} else {
		if err = d.Close(); err != nil {
			log.Fatal(err)
		}
	}
	return
}

// AssertDatabaseCount checks if a model's table has an expected amount of rows
func (suite *GormSuite) AssertDatabaseCount(t *testing.T, model interface{}, expected int64) {
	var count int64
	_ = suite.db.Model(&model).Count(&count)

	if assert.Equal(t, expected, count) {
		return
	}

	assert.Failf(t, "RequireDatabaseCount() expections not met", "model: %s, expected: %d, got: %d", reflect.ValueOf(model).Type().String(), expected, count)
}

// AssertDatabaseHas checks if a model is in the database
// TODO: Fix booleans
func (suite *GormSuite) AssertDatabaseHas(t *testing.T, model interface{}) {
	result := suite.db.Find(&model, model)

	if assert.Equal(t, int64(1), result.RowsAffected) {
		return
	}

	assert.Failf(t, "RequireDatabaseHas() expections not met", "model: %s, expected: %d, got: %d", reflect.ValueOf(model).Type().String(), 1, result.RowsAffected)
}

// AssertDatabaseMissing checks if a model is not in the database
func (suite *GormSuite) AssertDatabaseMissing(t *testing.T, model interface{}) {
	result := suite.db.Find(&model, model)

	if assert.Equal(t, int64(0), result.RowsAffected) {
		return
	}

	assert.Failf(t, "RequireDatabaseMissing() expections not met", "model: %s, expected: %d, got: %d", reflect.ValueOf(model).Type().String(), 0, result.RowsAffected)
}
