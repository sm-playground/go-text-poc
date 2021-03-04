package cacheclient

import (
	"errors"
	cnf "github.com/sm-playground/go-text-poc/config"
)

func GetCacheClient() (client CacheClient, err error) {
	config := cnf.GetInstance().Get()

	switch config.Cache.Dialect {
	case cnf.Redis:
		client = new(RedisClient)
	case cnf.Memcached:
		client = new(MemcachedClient)

	default:
		err = errors.New("undefined cache client")
	}

	if client != nil {
		client.InitCache()
	}

	return client, err
}
