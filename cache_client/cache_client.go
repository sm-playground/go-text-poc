package cacheclient

const NIL_VALUE_ERROR_MESSAGE = "cannot unmarshal nil value"

type CacheClient interface {
	// InitCache - initialize the cache
	InitCache()

	// Get - returns the cached value based on key
	Get(key string) (interface{}, error)

	// Set puts the value into the cache with the given key
	Set(key string, value interface{}) error

	// Delete removes the value with the given key from the cache
	Delete(key string) error

	// GetKeys returns all the keys from the cache matching the passed string pattern
	GetKeys(pattern string) ([]string, error)

	// Invalidates the data stored in cache for all keys matching the pattern
	Invalidate(pattern string) error
}
