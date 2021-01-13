package redisClient

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	c "github.com/sm-playground/go-text-poc/config"
	"log"
	"os"
	"strings"
)

var pool *redis.Pool

func InitCache(config c.Configurations) {
	// init redis connection pool
	initPool(config)

	// bootstrap some data to redis
	initStore()
}

// Returns redis connection.
func returnConnection(conn redis.Conn) {
	if err := conn.Close(); err != nil {
		log.Printf("ERROR!!! - failed to close redis connection")
	}
}

func Get(key string) (string, error) {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer returnConnection(conn)

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}

	return s, nil
}

func Set(key string, val string) error {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer returnConnection(conn)

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func Delete(key string) error {
	conn := pool.Get()
	defer returnConnection(conn)

	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Printf("ERROR: fail delete key %s, error %s", key, err.Error())
		return err
	}

	return nil

}

func initStore() {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer returnConnection(conn)

	macs := []string{"00000C  Cisco", "00000D  GOOGLE", "00000E  Fujitsu",
		"00000F  Next", "000010  Hughes"}
	for _, mac := range macs {
		pair := strings.Split(mac, "  ")
		err := Set(pair[0], pair[1])

		val, err := Get(pair[0])
		fmt.Println(val, err)

		err = Delete(pair[0])
	}
}

func initPool(config c.Configurations) {
	pool = &redis.Pool{
		MaxIdle:   config.Cache.DBCP.MaxIdle,
		MaxActive: config.Cache.DBCP.MaxActive,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(config.Cache.Network, fmt.Sprintf("%s:%d", config.Cache.IP, config.Cache.Port))
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}
