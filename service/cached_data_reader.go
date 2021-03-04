package service

import (
	"fmt"
	cache "github.com/sm-playground/go-text-poc/cache_client"
	"github.com/sm-playground/go-text-poc/common"
	m "github.com/sm-playground/go-text-poc/model"
)

// getCachedTextInfoList returns the cached value for the given cache key
func getCachedTextInfoList(cacheClient cache.CacheClient, cacheKey string) (textInfoList []m.TextInfo, err error) {

	var data interface{}

	if data, err = readFromCache(cacheClient, cacheKey); err != nil {
		return nil, err
	}

	var list m.TextInfoList
	l, e := list.Unmarshal(data)
	if e != nil && e.Error() != cache.NIL_VALUE_ERROR_MESSAGE {
		// Nothing to transform, no data found in the cache
		fmt.Printf("No data in the cache for key %s. Query from database", cacheKey)
	} else {
		if e == nil {
			// No error found and the data has been identified in the cache.
			return l.List, e
		}
	}

	return nil, nil
}

// getCachedTextInfo reads the data from the cache for the given key.
// Returns an interface
func getCachedTextInfo(cacheClient cache.CacheClient, cacheKey string) (textInfo interface{}, err error) {
	var data interface{}

	if data, err = readFromCache(cacheClient, cacheKey); err != nil {
		return textInfo, err
	}

	var ti m.TextInfo
	ti, e := ti.Unmarshal(data)
	if e != nil && e.Error() != cache.NIL_VALUE_ERROR_MESSAGE {
		// Nothing to transform, no data found in the cache
		fmt.Printf("No data in the cache for key %s. Query from database", cacheKey)
	} else {
		if e == nil {
			// No error found and the data has been identified in the cache.
			return ti, e
		}
	}

	return nil, nil
}

// getCachedTokenTextList returns the cached collection of TokenText objects
// for the given cache key
func getCachedTokenTextList(cacheClient cache.CacheClient, cacheKey string) (tokenTextList []m.TokenText, err error) {
	var data interface{}

	if data, err = readFromCache(cacheClient, cacheKey); err != nil {
		return nil, err
	}

	var list m.TokenTextList
	l, e := list.Unmarshal(data)
	if e != nil && e.Error() != cache.NIL_VALUE_ERROR_MESSAGE {
		// Nothing to transform, no data found in the cache
		fmt.Printf("No data in the cache for key %s. Query from database", cacheKey)
	} else {
		if e == nil {
			// No error found and the data has been identified in the cache.
			return l.List, e
		}
	}

	return nil, nil
}

// getReadDataCacheKey - returns the cacke key based on query input data
// getReadDataCacheKey returns the cache key for the read data request for the single token
func getReadDataCacheKey(queryInput SingleQueryInput, token string) string {
	locale := ""
	if queryInput.Language != "" {
		locale = queryInput.Language + "-"
	} else {
		locale = queryInput.DefaultLanguage + "-"
	}
	if queryInput.Country != "" {
		locale += queryInput.Country
	} else {
		locale += queryInput.DefaultCountry
	}

	key := common.GetServiceOwnerId() + ":" + locale + ":" + token

	if queryInput.TargetId != "" {
		key += ":" + queryInput.TargetId
	}

	return key
}
