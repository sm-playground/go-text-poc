package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	c "github.com/sm-playground/go-text-poc/common"
	m "github.com/sm-playground/go-text-poc/model"
	"net/http"
)

// DeleteTextInfo Deletes a single record from the text_info table
func DeleteTextInfo(r *http.Request, db *gorm.DB) (textInfo m.TextInfo, err error) {
	var requestStatus c.RequestStatus

	params := mux.Vars(r)
	var deletedRecordId = params["id"]
	if db.First(&textInfo, deletedRecordId).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", deletedRecordId))

	} else {
		err = nil
		db.Delete(&textInfo)
		requestStatus.Status = "success"
		requestStatus.Message = fmt.Sprintf("The record with id=%s was deleted", deletedRecordId)
	}

	return textInfo, err
}

// CreateTextInfo Creates a single record in the text_info table
func CreateTextInfo(textInfo m.TextInfo, db *gorm.DB) (err error) {
	err = nil
	if textInfo.Token != "" {
		db.Create(&textInfo)
	} else {
		err = errors.New("no error, try read")
	}

	return err
}

// UpdateTextInfo updates the record in the text_info table. ALL fields are updated
func UpdateTextInfo(r *http.Request, db *gorm.DB) (textInfo m.TextInfo, err error) {

	params := mux.Vars(r)
	var id = params["id"]
	if db.First(&textInfo, id).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", id))
		return textInfo, err
	} else {
		err = json.NewDecoder(r.Body).Decode(&textInfo)

		if err != nil {
			return textInfo, err
		}
		db.Save(&textInfo)

		return textInfo, nil
	}

}
