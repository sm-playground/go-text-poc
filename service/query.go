package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
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
func GetTextInfo(tokens []string, al string) (textInfoList []m.TextInfo) {
	db := d.GetConnection()

	if tokens != nil {

		if al == "" {
			al = "en-US"
		}

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

	return textInfoList

}

// GetSingleTextInfo returns the single record from the text_info table
//
// for the given id parameter in JSON format
func GetSingleTextInfo(params map[string]string) (textInfo m.TextInfo, err error) {

	db := d.GetConnection()

	if db.First(&textInfo, params["id"]).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", params["id"]))
	} else {
		err = nil
	}
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

	// sql := ""
	// var queryValues []interface{}

	for _, token := range payload.Tokens {

		data = append(data, resolveSingleToken(db, token, queryInput)...)

		// singleTokenSql, singleTokenValues := buildSingleTokenQuery(payload.TargetId, language, country, token)

		// @TODO for testing only. Remove this line
		// db.Raw(singleTokenSql, singleTokenValues...).Scan(&data)

		/*
			if i == 0 {
				sql = singleTokenSql
			} else {
				sql = sql + " UNION ALL " + singleTokenSql
			}
			queryValues = append(queryValues, singleTokenValues...)
		*/
	}

	// db.Raw(sql, queryValues...).Scan(&data)

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

// Returns the collection of TokenText objects.
//
// If the token represents a pattern more than a single record might be returned.
//
// If no records found the token is returned as a text.
func resolveSingleToken(db *gorm.DB, tokenInfo m.TokenPayload, queryInput SingleQueryInput) (data []m.TokenText) {
	var textInfo []m.TextInfo

	var query = ""
	var queryValues []interface{}
	token := tokenInfo.Token + "%"

	resolvedPlaceholder := resolvePlaceholders(tokenInfo)

	if queryInput.TargetId != "" {
		query, queryValues = getSubquery(queryInput, Customized, token, resolvedPlaceholder)
		query += " UNION "
	}

	s, v := getSubquery(queryInput, Localized, token, resolvedPlaceholder)
	query += s
	queryValues = append(queryValues, v...)

	if queryInput.Country != queryInput.DefaultCountry || queryInput.Language != queryInput.DefaultLanguage {
		// Add the localized fallout only if the queried locale is different than the default locale
		s, v = getSubquery(queryInput, FallbackLocale, token, resolvedPlaceholder)
		query += " UNION " + s
		queryValues = append(queryValues, v...)
	}

	s, v = getSubquery(queryInput, FallbackDefault, token, resolvedPlaceholder)
	query += " UNION " + s
	queryValues = append(queryValues, v...)

	query += " ORDER BY orderkey"

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

// Builds the complete SQL query requesting localized text information for a single token.
//
// Returns the SQL query with bind variables and the values in the separate array
func buildSingleTokenQuery(targetId string, language string, country string, tokenInfo m.TokenPayload) (query string, values []interface{}) {

	sqlSelect := "SELECT token, "
	sqlWhere := " WHERE target_id = ? AND token LIKE ?"
	sqlFrom := " FROM text_info"

	var whereValues []interface{}

	whereValues =
		append(whereValues, targetId, tokenInfo.Token+"%")

	if language != "" {
		sqlWhere += " AND language = ?"
		whereValues = append(whereValues, language)
	}
	if country != "" {
		sqlWhere += " AND country = ?"
		whereValues = append(whereValues, country)
	}

	resolvedPlaceholder := resolvePlaceholders(tokenInfo)

	sqlSelect += resolvedPlaceholder.SelectField
	if resolvedPlaceholder.PlaceholderValues != nil {
		values = append(values, resolvedPlaceholder.PlaceholderValues)
	}
	values = append(values, whereValues)

	query = sqlSelect + sqlFrom + sqlWhere

	return query, values
}
