package cacheclient

import (
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
)

type MemcachedClient struct {
}

func (rc *MemcachedClient) InitCache() {
	log.Error("Not implemented")
}

func (rc *MemcachedClient) Get(key string) (interface{}, error) {
	return "Not Implemented", errors.New(fmt.Sprintf("Not implemented - %s", key))
}

func (rc *MemcachedClient) Set(key string, value interface{}) error {
	return errors.New(fmt.Sprintf("Not implemented - %s, %v", key, value))
}

func (rc *MemcachedClient) Delete(key string) error {
	return errors.New(fmt.Sprintf("Not implemented - %s", key))
}

//noinspection GoUnusedParameter
func (rc *MemcachedClient) GetKeys(pattern string) ([]string, error) {
	log.Error("Not implemented")
	return nil, nil
}

//noinspection ALL
func (rc *MemcachedClient) Invalidate(pattern string) error {
	log.Error("Not implemented")
	return nil
}
