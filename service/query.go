package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	m "github.com/sm-playground/go-text-poc/model"
	"log"
	"net/http"
	"strings"
)

// GetTextInfo returns the list of records from the text_info database in JSON format
func GetTextInfo(r *http.Request, db *gorm.DB) (textInfoList []m.TextInfo) {

	tokens, ok := r.URL.Query()["token"]

	if ok && len(tokens[0]) > 0 {
		log.Printf("Url Param 'token' -> %s\n", tokens[0])

		al := r.Header.Get("Accept-Language")
		if al == "" {
			al = "en-US"
		}

		var values []interface{}
		query := ""

		als := strings.Split(al, ",")
		for _, locale := range als {
			// i - 0
			// locale en-CA
			l := strings.Split(locale, "-")

			if len(l) != 2 {
				// the locale is expected in the format of "en-US"
				// Ignore all other cases
				continue
			}

			if query == "" {
				query = "(locale = ? and country = ?)"
			} else {
				query += " or (locale = ? and country = ?)"

			}
			values = append(values, strings.TrimSpace(l[0]), strings.TrimSpace(l[1]))
		}

		query = "(" + query + ")"

		query += " and token = ?"
		values = append(values, tokens[0])

		db.Where(query, values...).Find(&textInfoList)

	} else {
		db.Find(&textInfoList)
	}

	return textInfoList

}

// GetSingleTextInfo returns the single record from the text_info table
//
// for the given id parameter in JSON format
func GetSingleTextInfo(r *http.Request, db *gorm.DB) (textInfo m.TextInfo, err error) {
	params := mux.Vars(r)

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
func ReadData(r *http.Request, db *gorm.DB) (data []m.TokenText, err error) {
	var payload m.TextInfoPayload

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	sql := ""
	var queryValues []interface{}
	for i, token := range payload.Tokens {

		singleTokenSql, singleTokenValues := buildSingleTokenQuery(payload, token)

		db.Raw(singleTokenSql, singleTokenValues...).Scan(&data)

		if i == 0 {
			sql = singleTokenSql
		} else {
			sql = sql + " UNION ALL " + singleTokenSql
		}
		queryValues = append(queryValues, singleTokenValues...)

	}

	db.Raw(sql, queryValues...).Scan(&data)

	return data, nil

}

// Builds the complete SQL query requesting localized text information for a single token.
//
// Returns the SQL query with bind variables and the values in the separate array
func buildSingleTokenQuery(payload m.TextInfoPayload, tokenInfo m.TokenPayload) (query string, values []interface{}) {

	sqlSelect := "SELECT token, "
	sqlWhere := " WHERE target_id = ? AND locale = ? AND country = ? AND token LIKE ?"
	sqlFrom := " FROM text_info"

	var whereValues []interface{}

	whereValues =
		append(whereValues, payload.TargetId, payload.Locale, payload.Country, tokenInfo.Token+"%")

	if len(tokenInfo.Placeholders) == 0 {
		// No placeholders.
		sqlSelect += "text"
		values = whereValues
	} else {
		// There is a collection of placeholders in the request
		sqlSelectPlaceholder := ""
		var placeholderValues []interface{}
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
		sqlSelect += sqlSelectPlaceholder + " as text"

		values = append(placeholderValues, whereValues...)
	}

	query = sqlSelect + sqlFrom + sqlWhere

	return query, values
}
