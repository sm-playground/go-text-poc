package cacheclient

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	c "github.com/sm-playground/go-text-poc/config"
	"log"
	"os"
)

const CACHE_EXPIRATION_IN_MINUTES = 10 * 60

type RedisClient struct {
}

var pool *redis.Pool

func (rc *RedisClient) InitCache() {
	// init redis connection pool
	initPool()

}

// Get returns the cached value based on key
//
// The caller is responsible for unmarshalling the data into proper object
func (rc *RedisClient) Get(key string) (interface{}, error) {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer returnConnection(conn)

	val, err := conn.Do("GET", key)

	return val, err
}

// Set puts the value into the cache with the given key
func (rc *RedisClient) Set(key string, value interface{}) (err error) {
	conn := pool.Get()
	defer returnConnection(conn)

	if j, err := json.Marshal(value); err == nil {
		_, err = conn.Do("SET", key, j, "EX", CACHE_EXPIRATION_IN_MINUTES)
	}

	return err
}

// Delete removes the value with the given key from the cache
func (rc *RedisClient) Delete(key string) error {
	conn := pool.Get()
	defer returnConnection(conn)

	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Printf("ERROR: fail delete key %s, error %s", key, err.Error())
		return err
	}

	return nil
}

func (rc *RedisClient) GetKeys(pattern string) ([]string, error) {
	conn := pool.Get()
	defer returnConnection(conn)

	iter := 0
	var keys []string
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	fmt.Printf("keys %+v", keys)

	return keys, nil

}

func (rc *RedisClient) Invalidate(pattern string) error {
	conn := pool.Get()
	defer returnConnection(conn)

	if keys, err := rc.GetKeys(pattern); err == nil {
		for _, key := range keys {
			err = rc.Delete(key)
		}
	}

	return nil
}

// Returns redis connection.
func returnConnection(conn redis.Conn) {
	if err := conn.Close(); err != nil {
		log.Printf("ERROR!!! - failed to close redis connection")
	}
}

// initPool Initializes the redis connection pool based on config parameters
func initPool() {
	config := c.GetInstance().Get()

	pool = &redis.Pool{
		MaxIdle:   config.Cache.Pool.MaxIdle,
		MaxActive: config.Cache.Pool.MaxActive,
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
