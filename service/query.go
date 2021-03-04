package service

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	cache "github.com/sm-playground/go-text-poc/cache_client"
	"github.com/sm-playground/go-text-poc/common"
	c "github.com/sm-playground/go-text-poc/config"
	d "github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"strings"
)

type SubqueryOrderKey string

const (
	Customized      SubqueryOrderKey = "a"
	Localized       SubqueryOrderKey = "b"
	FallbackLocale  SubqueryOrderKey = "c"
	FallbackDefault SubqueryOrderKey = "d"
)

type SingleQueryInput struct {
	TargetId             string
	SourceId             string
	Country              string
	Language             string
	DefaultCountry       string
	DefaultLanguage      string
	ServiceOwnerSourceId string
}

type ResolvedPlaceholder struct {
	SelectField       string
	PlaceholderValues []interface{}
}

// GetTextInfo returns the list of records from the text_info database in JSON format
func GetTextInfo(tokens []string, al string) (textInfoList []m.TextInfo, err error) {

	var cacheClient cache.CacheClient
	if cacheClient, err = cache.GetCacheClient(); err != nil {
		return nil, err
	}

	db := d.GetConnection()

	if al == "" {
		al = c.GetInstance().Get().DefaultLocale
	}

	// First try to read the data from the cache
	cacheKey := common.GetServiceOwnerId() + ":" + al + ":" + strings.Join(tokens, "-")
	textInfoList, err = getCachedTextInfoList(cacheClient, cacheKey)
	if textInfoList != nil && err == nil {
		// Found the data in the cache
		fmt.Printf("Found data in the cache for key - %s\n", cacheKey)
		return textInfoList, err
	}

	// The data is not found in the cache. Query it and return
	if tokens != nil {
		var values []interface{}
		query := ""

		als := strings.Split(al, ",")
		for _, locale := range als {
			// locale en-CA
			l := strings.Split(locale, "-")

			if len(l) != 2 {
				// the locale is expected in the format of "en-US"
				// Ignore all other cases
				continue
			}

			language := strings.TrimSpace(l[0])
			country := strings.TrimSpace(l[1])

			// Query all the records with matching language and country
			// OR only language if the country IS NULL
			if query == "" {
				query = "((language = ? and country = ?) or (language = ? and country = ''))"
			} else {
				query += " or ((language = ? and country = ?) or (language = ? and country = ''))"

			}
			values = append(values, language, country, language)
		}

		query = "(" + query + ")"

		query += " and token like ?"
		values = append(values, tokens[0]+"%")

		db.Where(query, values...).Find(&textInfoList)

	} else {
		db.Find(&textInfoList)
	}

	// put the list into the wrapper object and cache it
	var list m.TextInfoList
	list.SetList(textInfoList)
	fmt.Printf("\nput the list into the wrapper object and cache it - %s\n", cacheKey)
	err = cacheClient.Set(cacheKey, list)

	return textInfoList, err

}

// GetSingleTextInfo returns the single record from the text_info table
//
// for the given id parameter in JSON format
func GetSingleTextInfo(params map[string]string) (textInfo m.TextInfo, err error) {

	var cacheClient cache.CacheClient
	if cacheClient, err = cache.GetCacheClient(); err != nil {
		return textInfo, err
	}

	cacheKey := common.GetServiceOwnerId() + ":TEXT_INFO:" + params["id"]

	tiCached, er := getCachedTextInfo(cacheClient, cacheKey)
	if tiCached != nil && er == nil {
		fmt.Printf("Found data in the cache for key - %s\n", cacheKey)
		return tiCached.(m.TextInfo), err
	}

	db := d.GetConnection()

	if db.First(&textInfo, params["id"]).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", params["id"]))
	} else {
		err = nil
	}

	fmt.Printf("cache the object - %s\n", cacheKey)
	err = cacheClient.Set(cacheKey, textInfo)

	return textInfo, err
}

// ReadData returns the localized text information according to request data provided in the payload.
//
// The function is a handler for the POST requests. The returned data is filtered by country, locale, and
// target (company Id) and capable of processing multiple token with and without the data for the placeholders.
func ReadData(payload m.TextInfoPayload) (data []m.TokenText, err error) {
	db := d.GetConnection()

	config := c.GetInstance().Get()

	var language, country string

	language, country, err = resolveLocale(payload)

	if err != nil {
		// return an error if the locale cannot be resolved
		return data, err
	}

	queryInput := SingleQueryInput{
		ServiceOwnerSourceId: config.ServiceOwnerSourceId,
		DefaultCountry:       config.DefaultCountry,
		DefaultLanguage:      config.DefaultLanguage,
		Country:              country,
		Language:             language,
		TargetId:             payload.TargetId,
		SourceId:             payload.SourceId}

	for _, token := range payload.Tokens {

		data = append(data, resolveSingleToken(db, token, queryInput)...)

	}

	return data, nil

}

// getSubquery - returns the query and the bind parameters based on the subqueryOrderKey discriminator
// identifying the type of the query requested by the caller.
//
// The supported query types are: customized, localized, fallout localized, and fallout default
func getSubquery(queryInput SingleQueryInput, subqueryOrderKey SubqueryOrderKey, token string, resolvedPlaceholder ResolvedPlaceholder) (subquery string, subqueryValues []interface{}) {

	querySelect := "SELECT ? as orderkey, target_id, is_fallback , token, " + resolvedPlaceholder.SelectField + " "

	queryFrom := "FROM text_info "

	queryWhere := "WHERE token like ? "

	// bind order key
	subqueryValues = append(subqueryValues, subqueryOrderKey)
	// bind placeholder values
	if resolvedPlaceholder.PlaceholderValues != nil {
		subqueryValues = append(subqueryValues, resolvedPlaceholder.PlaceholderValues...)
	}
	subqueryValues = append(subqueryValues, token)

	switch subqueryOrderKey {
	case Customized:
		// Include language and country unconditionally
		queryWhere += " and language = ? and country = ? "
		subqueryValues = append(subqueryValues, queryInput.Language, queryInput.Country)
		if queryInput.TargetId != "" {
			// include the target Id in the request if provided
			queryWhere += " and target_id = ? "
			subqueryValues = append(subqueryValues, queryInput.TargetId)
		}
		break

	case Localized:
		// Include language and country in the query for all NON-FALLBACK records
		queryWhere += " and language = ? and country = ?  and source_id = ? and NOT is_fallback"
		subqueryValues = append(subqueryValues, queryInput.Language, queryInput.Country, queryInput.ServiceOwnerSourceId)
		break

	case FallbackLocale:
		queryWhere += " and language = ? and source_id = ? and is_fallback"
		subqueryValues = append(subqueryValues, queryInput.Language, queryInput.ServiceOwnerSourceId)
		break

	case FallbackDefault:
		queryWhere += " and language = ? and source_id = ? and is_fallback"
		subqueryValues = append(subqueryValues, queryInput.DefaultLanguage, queryInput.ServiceOwnerSourceId)
		break

	default:
		return "", nil
	}

	query := querySelect + queryFrom + queryWhere

	return query, subqueryValues
}

// readFromCache - Tries to read from the cache for the given key
func readFromCache(cacheClient cache.CacheClient, key string) (data interface{}, err error) {

	if data, err = cacheClient.Get(key); err != nil && err != redis.ErrNil {
		// Got error while reading data from the cache
		return nil, err
	}

	return data, nil
}

// Returns the collection of TokenText objects.
//
// If the token represents a pattern more than a single record might be returned.
//
// If no records found the token is returned as a text.
func resolveSingleToken(db *gorm.DB, tokenInfo m.TokenPayload, queryInput SingleQueryInput) (data []m.TokenText) {

	token := tokenInfo.Token + "%"

	// First check the data in the cache
	cacheKey := getReadDataCacheKey(queryInput, token)
	var cacheClient cache.CacheClient
	var err error
	if cacheClient, err = cache.GetCacheClient(); err != nil {
		return nil
	}
	tokenTextList, err := getCachedTokenTextList(cacheClient, cacheKey)

	if tokenTextList != nil && err == nil {
		// Found the data in the cache
		fmt.Printf("Found data in the cache for key - %s\n", cacheKey)
		return tokenTextList
	}

	// No data in the cache for that key. Query from the database
	var query = ""
	var queryValues []interface{}
	resolvedPlaceholder := resolvePlaceholders(tokenInfo)

	if queryInput.TargetId != "" {
		// Include customized records if target id is provided
		query, queryValues = getSubquery(queryInput, Customized, token, resolvedPlaceholder)
		query += " UNION "
	}

	// Include localized records
	s, v := getSubquery(queryInput, Localized, token, resolvedPlaceholder)
	query += s
	queryValues = append(queryValues, v...)

	if queryInput.Country != queryInput.DefaultCountry || queryInput.Language != queryInput.DefaultLanguage {
		// Add the localized fallout only if the queried locale is different than the default locale
		s, v = getSubquery(queryInput, FallbackLocale, token, resolvedPlaceholder)
		query += " UNION " + s
		queryValues = append(queryValues, v...)
	}

	// Include default fallback records
	s, v = getSubquery(queryInput, FallbackDefault, token, resolvedPlaceholder)
	query += " UNION " + s
	queryValues = append(queryValues, v...)

	query += " ORDER BY orderkey"

	var textInfo []m.TextInfo
	db.Raw(query, queryValues...).Scan(&textInfo)

	if len(textInfo) == 0 {
		// No record is found for the given token. Return the token
		data = append(data, m.TokenText{Text: tokenInfo.Token, Token: tokenInfo.Token})
	} else {
		var pt = ""
		for _, record := range textInfo {
			// fmt.Printf("\n%s", pt)
			if pt != record.Token {
				pt = record.Token
				// If the token has more than one record they are ordered as customized - localized - fallback
				// take the first one and ignore the rest
				// fmt.Printf("\n%+v", record)
				data = append(data, m.TokenText{Text: record.Text, Token: record.Token})
			} else {
				continue
			}
		}
	}

	var list m.TokenTextList
	list.SetList(data)
	fmt.Printf("\nput the list into the wrapper object and cache it - %s\n", cacheKey)
	err = cacheClient.Set(cacheKey, list)

	return data
}

// resolvePlaceholders returns the ResolvedPlaceholder structure with queried field (optionally, with formula)
// and the bind variables
func resolvePlaceholders(tokenInfo m.TokenPayload) (resolvedPlaceholder ResolvedPlaceholder) {

	var selectField string
	var placeholderValues []interface{}

	if len(tokenInfo.Placeholders) == 0 {
		// No placeholders.
		selectField = "text"
	} else {
		// There is a collection of placeholders in the request
		sqlSelectPlaceholder := ""
		// Iterate over the placeholders:
		// Concatenate the final sql request with bind variables
		// Keep collecting the values that will be applied to the query
		for i, ph := range tokenInfo.Placeholders {
			if i == 0 {
				sqlSelectPlaceholder = "replace (text, ?, ?)"
				placeholderValues = append(placeholderValues, "{"+ph.Name+"}", ph.Value)
			} else {
				sqlSelectPlaceholder = "replace(" + sqlSelectPlaceholder + ", ?, ?)"
				placeholderValues = append(placeholderValues, "{"+ph.Name+"}", ph.Value)
			}
		}
		selectField = sqlSelectPlaceholder + " as text"
	}

	return ResolvedPlaceholder{selectField, placeholderValues}
}

// The method takes the payload from the requests and reads the provided locale. In case of missing locale information
// in the request the defaults from the config is returned
func resolveLocale(payload m.TextInfoPayload) (language string, country string, err error) {
	if strings.TrimSpace(payload.Locale) == "" {
		// No locale is provided. Take the language and the country from the payload
		language = payload.Language
		country = payload.Country
	} else {
		// If the locale is provided it takes precedence and country and language are retrieved from locale
		l := strings.Split(payload.Locale, "-")

		if len(l) != 2 {
			// the locale is expected in the format of "en-US"
			// Ignore all other cases
			return language, country, errors.New(fmt.Sprintf("incorrect format for locale: %s", payload.Locale))
		}
		language = strings.TrimSpace(l[0])
		country = strings.TrimSpace(l[1])
	}

	config := c.GetInstance().Get()

	if strings.TrimSpace(language) == "" {
		language = config.DefaultLanguage
	}
	if strings.TrimSpace(country) == "" {
		country = config.DefaultCountry
	}

	return language, country, nil
}
