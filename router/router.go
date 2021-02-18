package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	cm "github.com/sm-playground/go-text-poc/common"
	c "github.com/sm-playground/go-text-poc/config"
	m "github.com/sm-playground/go-text-poc/model"
	s "github.com/sm-playground/go-text-poc/service"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

type PostRequestHandler struct{}

// ServeHTTP is a request interceptor implementing the following logic.
//
// 1. Clones the request because the body of the request can be retrieved just once.
//
// 2. Tries to read the content of the body into the TextInfo structure
//
// 3. Error at #2 is an indication of request with the method other than POST. Redirect forward
//
// 4. In case of no error and the body contains the parameter Token call the CreateTextInfo method
//
// 5. If no error and no Token redirect to the ReadTextInfo method with the cloned request as a parameter
func (*PostRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	t := time.Now()

	/*
		r2 := r.Clone(r.Context())
		*r2 = *r

		var b bytes.Buffer
		_, err := b.ReadFrom(r.Body)
		if err != nil {
			log.Printf("ERROR!!! - failed reading the request body")
			return
		}

		r.Body = ioutil.NopCloser(&b)
		r2.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

		var textInfo model.TextInfo

		if json.NewDecoder(r.Body).Decode(&textInfo) != nil {
			// Got request with method other than POST
			fmt.Printf("Received request with method - %s. Move forward to serve the request", r.Method)
			next.ServeHTTP(w, r2)
		} else {
			if textInfo.Token != "" {
				fmt.Printf("POST request received with textInfo valid payload. Redirect to CreateTextInfo")
				CreateTextInfo(w, textInfo)
			} else {
				// Received the request with the POST method to retrieve the data according to payload.
				// Move forward to ReadTextInfo
				fmt.Printf("Received request with method - %s. Move forward to ReadTextInfo", r.Method)
				next.ServeHTTP(w, r2)
			}
		}
	*/

	// Forward to the request handler
	next.ServeHTTP(w, r)

	fmt.Printf("Execution time: %s \n", time.Now().Sub(t).String())
}

// Initializes all supported routers
func InitRouter() {
	config := c.GetInstance().Get()

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()

	// regular router
	router.HandleFunc("/textInfo", GetTextInfo).Methods("GET")

	router.HandleFunc("/textInfo/{id:[0-9]+}", GetSingleTextInfo).Methods("GET")

	router.HandleFunc("/textInfo/{id:[0-9]+}", DeleteTextInfo).Methods("DELETE")

	router.HandleFunc("/textInfo/{id:[0-9]+}", OverwriteTextInfo).Methods("PUT")

	router.HandleFunc("/textInfo/{id:[0-9]+}", UpdateTextInfo).Methods("PATCH")

	// router.HandleFunc("/textInfo", ReadTextInfo).Methods("POST")

	router.HandleFunc("/textInfo", CreateTextInfo).Methods("POST")

	router.HandleFunc("/textInfo/query", ReadTextInfo).Methods("POST")

	// batch requests
	router.HandleFunc("/textInfo/batch", BatchCreate).Methods("POST")

	router.HandleFunc("/textInfo/batch", BatchUpdate).Methods("PUT")

	router.HandleFunc("/textInfo/batch", BatchDelete).Methods("DELETE")

	n := negroni.New()
	n.Use(&PostRequestHandler{})
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port), n))
}

// GetTextInfo GET method handler.
// Wraps a call to the query service method returning the complete list of the TextInfo records.
// Returns the list in the JSON format in the response
func GetTextInfo(w http.ResponseWriter, r *http.Request) {

	var textInfoList = s.GetTextInfo(
		r.URL.Query()["token"],
		r.Header.Get("Accept-Language"))

	if json.NewEncoder(w).Encode(&textInfoList) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// GetSingleTextInfo GET method handler
// Wraps a call to the query service returning a single record from the text_info table for the given id.
//
// Returns the record in the JSON format in the response
func GetSingleTextInfo(w http.ResponseWriter, r *http.Request) {

	textInfo, err := s.GetSingleTextInfo(mux.Vars(r))
	if err != nil {
		var requestStatus cm.RequestStatus
		requestStatus.Status = "failed"
		requestStatus.Message = err.Error()
		if json.NewEncoder(w).Encode(&requestStatus) != nil {
			log.Printf("ERROR!!! - failed encoding the query response")
		}
	} else {
		if json.NewEncoder(w).Encode(&textInfo) != nil {
			log.Printf("ERROR!!! - failed encoding the query response")
		}
	}
}

// DeleteTextInfo DELETE method handler
// Wraps a call to the processing service to delete a single record from the text_info table
func DeleteTextInfo(w http.ResponseWriter, r *http.Request) {

	var requestStatus cm.RequestStatus
	textInfo, err := s.DeleteTextInfo(mux.Vars(r))
	if err != nil {
		requestStatus.Status = "failed"
		requestStatus.Message = err.Error()
	} else {
		requestStatus.Status = "success"
		requestStatus.Message = fmt.Sprintf("The record with id=%d was deleted", textInfo.Id)
	}

	if json.NewEncoder(w).Encode(&requestStatus) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// UpdateTextInfo updates the record in the text_info table. Only specified fields are updated
func UpdateTextInfo(w http.ResponseWriter, r *http.Request) {
	textInfo, err := s.UpdateTextInfo(mux.Vars(r), r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&textInfo) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// OverwriteTextInfo overwrites the record in the text_info table. ALL fields are updated
func OverwriteTextInfo(w http.ResponseWriter, r *http.Request) {

	textInfo, err := s.OverwriteTextInfo(mux.Vars(r), r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&textInfo) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// CreateTextInfo POST method handler
//
// Wraps a call to the processing service to create a single record in the text_info table
func CreateTextInfo(w http.ResponseWriter, r *http.Request) {

	var ti m.TextInfo
	er := json.NewDecoder(r.Body).Decode(&ti)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	ti, er = s.CreateTextInfo(ti)

	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", ti)

	if json.NewEncoder(w).Encode(&ti) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// ReadTextInfo POST method handler
//
// Wraps a call to the query service to read the records from the database for the
// query payload provided in the body of the request
func ReadTextInfo(w http.ResponseWriter, r *http.Request) {

	var payload m.TextInfoPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tokenTextList []m.TokenText

	tokenTextList, err = s.ReadData(payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&tokenTextList) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// - - - - - - - Batch processes - - - - - - - -

// BatchCreate POST method handler
//
// Wraps a call to the query service to create multiple records in the text_info table
// for the collection of objects in the JSON format
func BatchCreate(w http.ResponseWriter, r *http.Request) {
	response, err := s.BatchCreate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&response) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// BatchUpdate PUT method handler
//
// Wraps a call to the query service to update (overwrite) multiple records in the text_info table
// for the collection of objects in the JSON format
func BatchUpdate(w http.ResponseWriter, r *http.Request) {
	response, err := s.BatchUpdate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&response) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}

// BatchDelete DELETE method handler
//
// Wraps a call to the query service to delete multiple records in the text_info table
// for the collection of objects in the JSON format
func BatchDelete(w http.ResponseWriter, r *http.Request) {
	response, err := s.BatchDelete(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(&response) != nil {
		log.Printf("ERROR!!! - failed encoding the query response")
	}
}
