package database

import (
	"fmt"

	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/tokopedia/sqlt"
)

var DBPool struct {
	MainDB *sqlt.DB
}

var RedisPool struct {
	RedisDev *redis.Pool
}

func InitDB() {
	databaseString := "postgres://kunyit_user:Ti8yN7uk65@devel-postgre.tkpd/tokopedia-order?sslmode=disable"
	databaseCon := databaseString + ";" + databaseString
	db, err := sqlt.Open("postgres", databaseCon)
	if err != nil {
		fmt.Println(err)
	}

	DBPool.MainDB = db
}

func InitRedis() {
	address := "localhost:6379"
	RedisPool.RedisDev = generatePool(address)
}

func generatePool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
