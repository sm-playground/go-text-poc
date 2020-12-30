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
)

type TextInfo struct {
	// gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Token      string `gorm:"type:varchar(100)"`
	Text       string `gorm:"type:varchar(2000)"`
	Noun       string
	Function   string
	Action     string `gorm:"size:30"`
	Module     string
	Country    string
	Locale     string
	SourceType string
	SourceId   string
	TargerId   string
	ReadOnly   bool
}

type RequestStatus struct {
	Status  string
	Message string
}

var db *gorm.DB
var err error

var (
	textInfo = []TextInfo{
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Record type", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Record type", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.ENTITY", Text: "Customer entity", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.ENTITY", Text: "Customer entity", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Record type", Action: "View", Country: "CA", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Type d'enregistrement", Action: "View", Country: "CA", Locale: "fr", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "View", Country: "CA", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Nom du client", Action: "View", Country: "CA", Locale: "fr", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "ომხმარებლის სახელი", Action: "View", Country: "GE", Locale: "ge", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Հաճախորդի անունը", Action: "View", Country: "AM", Locale: "am", ReadOnly: true},
		// {Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Գրանցման տեսակը", Action: "View", Country: "AM", Locale: "am", ReadOnly: true},
	}
)

func main() {

	// Initialize the database and populate with the sample data
	db = initDatabase()
	defer db.Close()

	// Initialize redis connection pool
	redisClient.InitRedis()

	redisClient.Set("hello", "hello world")

	var hello, err = redisClient.Get("hello")

	fmt.Println(hello, err)

	router := mux.NewRouter()

	router.HandleFunc("/textInfo", GetTextInfo).Methods("GET")
	router.HandleFunc("/textInfo/{id}", GetSingleTextInfo).Methods("GET")
	router.HandleFunc("/textInfo", CreateTextInfo).Methods("POST")
	router.HandleFunc("/textInfo/{id}", UpdateTextInfo).Methods("PUT")
	router.HandleFunc("/textInfo/{id}", DeleteTextInfo).Methods("DELETE")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

// initDatabase initializes the database
// - connect to postgres
// - runs AutoMigrate to create the database
// - populate the text_info table with the localized records list
func initDatabase() (db *gorm.DB) {

	db, err = gorm.Open("postgres", "port=54320 user=postgres dbname=ia_text sslmode=disable password=postgres")

	if err != nil {
		panic("failed to connect database")
	}

	// set the log mode to see the queries executed by the gorm
	db.LogMode(true)

	// instructs gorm to create the tables with the singular names
	db.SingularTable(true)

	// create tables based on specified structs
	db.AutoMigrate(&TextInfo{})

	for index := range textInfo {
		db.Create(&textInfo[index])
	}

	return db
}

// GetTextInfo returns the list of records from the text_info database in JSON format
func GetTextInfo(w http.ResponseWriter, r *http.Request) {

	var textInfoList []TextInfo

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

	json.NewEncoder(w).Encode(&textInfoList)
}

// GetSingleTextInfo returns the single record from the text_info table
// for the given id parameter in JSON format
func GetSingleTextInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var textInfo TextInfo
	if db.First(&textInfo, params["id"]).RecordNotFound() {
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
	var textInfo TextInfo
	if db.First(&textInfo, deletedRecordId).RecordNotFound() {
		requestStatus.Status = "failed"
		requestStatus.Message = fmt.Sprintf("The record with id=%s is not found", deletedRecordId)
	} else {
		db.Delete(&textInfo)
		requestStatus.Status = "success"
		requestStatus.Message = fmt.Sprintf("The record with id=%s was deleted", deletedRecordId)
	}

	json.NewEncoder(w).Encode(&requestStatus)
}

// CreateTextInfo Creates a single record in the text_info table
func CreateTextInfo(w http.ResponseWriter, r *http.Request) {
	var textInfo TextInfo
	err = json.NewDecoder(r.Body).Decode(&textInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Create(&textInfo)

	json.NewEncoder(w).Encode(&textInfo)
}

// UpdateTextInfo updates the record in the text_info table. ALL fields are updated
func UpdateTextInfo(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var id = params["id"]
	var textInfo TextInfo
	if db.First(&textInfo, id).RecordNotFound() {
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
		db.Save(&textInfo)
		json.NewEncoder(w).Encode(&textInfo)
	}

}
