package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/sm-playground/go-text-poc/redisClient"
	"log"
	"net/http"
	"strings"

	c "github.com/sm-playground/go-text-poc/config"
	"github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
)

type RequestStatus struct {
	Status  string
	Message string
}

var dbClient *gorm.DB
var err error
var config c.Configurations

func main() {

	// Read configuration parameters
	config = c.LoadConfig()

	// Initialize the database and populate with the sample data
	dbClient = db.InitDatabase(config)

	defer dbClient.Close()

	// Initialize redis connection pool
	redisClient.InitCache(config)

	redisClient.Set("hello", "hello world")

	var hello, err = redisClient.Get("hello")

	fmt.Println(hello, err)

	// router := mux.NewRouter()
	// router.HandleFunc("/textInfo", GetTextInfo).Methods("GET")

	subrouter := mux.NewRouter().PathPrefix("/v1").Subrouter()
	subrouter.HandleFunc("/textInfo", GetTextInfo).Methods("GET")

	subrouter.HandleFunc("/textInfo/{id}", GetSingleTextInfo).Methods("GET")
	subrouter.HandleFunc("/textInfo", CreateTextInfo).Methods("POST")
	subrouter.HandleFunc("/textInfo/{id}", UpdateTextInfo).Methods("PUT")
	subrouter.HandleFunc("/textInfo/{id}", DeleteTextInfo).Methods("DELETE")

	handler := cors.Default().Handler(subrouter)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port), handler))
}

// GetTextInfo returns the list of records from the text_info database in JSON format
func GetTextInfo(w http.ResponseWriter, r *http.Request) {

	var textInfoList []m.TextInfo

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

		dbClient.Where(query, values...).Find(&textInfoList)

	} else {
		dbClient.Find(&textInfoList)
	}

	json.NewEncoder(w).Encode(&textInfoList)
}

// GetSingleTextInfo returns the single record from the text_info table
// for the given id parameter in JSON format
func GetSingleTextInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var textInfo m.TextInfo
	if dbClient.First(&textInfo, params["id"]).RecordNotFound() {
		var requestStatus RequestStatus
		requestStatus.Status = "failed"
		requestStatus.Message = fmt.Sprintf("The record with id=%s is not found", params["id"])
		json.NewEncoder(w).Encode(&requestStatus)
	} else {
		json.NewEncoder(w).Encode(&textInfo)
	}
}

// DeleteTextInfo Deletes a single record from the text_info table
func DeleteTextInfo(w http.ResponseWriter, r *http.Request) {
	var requestStatus RequestStatus

	params := mux.Vars(r)
	var deletedRecordId string = params["id"]
	var textInfo m.TextInfo
	if dbClient.First(&textInfo, deletedRecordId).RecordNotFound() {
		requestStatus.Status = "failed"
		requestStatus.Message = fmt.Sprintf("The record with id=%s is not found", deletedRecordId)
	} else {
		dbClient.Delete(&textInfo)
		requestStatus.Status = "success"
		requestStatus.Message = fmt.Sprintf("The record with id=%s was deleted", deletedRecordId)
	}

	json.NewEncoder(w).Encode(&requestStatus)
}

// CreateTextInfo Creates a single record in the text_info table
func CreateTextInfo(w http.ResponseWriter, r *http.Request) {
	var textInfo m.TextInfo
	err = json.NewDecoder(r.Body).Decode(&textInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbClient.Create(&textInfo)

	json.NewEncoder(w).Encode(&textInfo)
}

// UpdateTextInfo updates the record in the text_info table. ALL fields are updated
func UpdateTextInfo(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var id = params["id"]
	var textInfo m.TextInfo
	if dbClient.First(&textInfo, id).RecordNotFound() {
		var requestStatus RequestStatus
		requestStatus.Status = "failed"
		requestStatus.Message = fmt.Sprintf("The record with id=%s is not found", id)
		json.NewEncoder(w).Encode(&requestStatus)
	} else {
		err = json.NewDecoder(r.Body).Decode(&textInfo)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dbClient.Save(&textInfo)
		json.NewEncoder(w).Encode(&textInfo)
	}

}
