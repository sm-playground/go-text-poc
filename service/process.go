package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	c "github.com/sm-playground/go-text-poc/common"
	cnf "github.com/sm-playground/go-text-poc/config"
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
func CreateTextInfo(r *http.Request, db *gorm.DB, config cnf.Configurations) (textInfo m.TextInfo, err error) {

	err = json.NewDecoder(r.Body).Decode(&textInfo)
	if err != nil {
		return textInfo, err
	}

	if textInfo.TargetId == "" {
		// No target Id is specified so this record is applicable to all targets
		// Hence, the source shall be set to the service owner source Id
		textInfo.SourceId = config.ServiceOwnerSourceId
	}

	var count int64
	db.Model(&m.TextInfo{}).
		Where("token = ? AND language = ? AND is_fallback", textInfo.Token, textInfo.Language).
		Count(&count)
	if count == 0 {
		// No fallback is found for the token / locale.
		// Make it a fallback record and set the service owner source Id
		textInfo.IsFallback = true
	} else {
		if textInfo.IsFallback {
			// error condition. the default record is already in the database and
			// the user tries to insert another default
			return textInfo, errors.New("the fallback record is already in the database and tries to insert another one")
		}
	}

	if textInfo.Token != "" {
		db.Create(&textInfo)
	} else {
		err = errors.New("no token value found")
	}

	return textInfo, err
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
