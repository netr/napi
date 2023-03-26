package sweets

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

type RedisSuite struct {
	db *redis.Client
}

// NewRedisSuite is used to instantiate a new redis test suite. Typically called in SetupSuite(). Uses miniredis (https://github.com/alicebob/miniredis/v2) under the hood.
func (suite *RedisSuite) NewRedisSuite(t *testing.T) error {
	s := miniredis.RunT(t)
	db := redis.NewClient(&redis.Options{
		Addr:         s.Addr(),
		DB:           0, // use default DB
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	})

	suite.db = db
	return nil
}

// DB is a helper function to retrieve the underlying redis.Client
func (suite *RedisSuite) DB() *redis.Client {
	return suite.db
}
