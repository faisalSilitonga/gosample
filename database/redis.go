package database

import "net/http"
import "fmt"
import "github.com/garyburd/redigo/redis"

func RedisSetHandler(w http.ResponseWriter, r *http.Request) {
	redisConn := RedisPool.RedisDev.Get()
	defer redisConn.Close()

	key := r.FormValue("key")
	value := r.FormValue("val")

	_, err := redisConn.Do("SET", key, value)
	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "success set %s to %s", value, key)
	return
}

func RedisGetHandler(w http.ResponseWriter, r *http.Request) {
	redisConn := RedisPool.RedisDev.Get()
	defer redisConn.Close()

	key := r.FormValue("key")

	val, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "got %s", val)
	return
}
