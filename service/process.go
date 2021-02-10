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

	err = createTextInfo(&textInfo, db, config)

	return textInfo, err
}

// UpdateTextInfo updates the record in the text_info table. Only the fields specified in the request are updated
func UpdateTextInfo(r *http.Request, db *gorm.DB) (textInfo m.TextInfo, err error) {

	params := mux.Vars(r)
	var id = params["id"]
	if db.First(&textInfo, id).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", id))
		return textInfo, err
	} else {
		var newPyload m.TextInfo

		err = json.NewDecoder(r.Body).Decode(&newPyload)
		if err != nil {
			return textInfo, err
		}

		textInfo.Overwrite(newPyload)

		db.Save(&textInfo)

		return textInfo, nil
	}
}

// OverwriteTextInfo updates the record in the text_info table. ALL fields are updated
func OverwriteTextInfo(r *http.Request, db *gorm.DB) (textInfo m.TextInfo, err error) {

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

// createTextInfo creates a single record in the text_info table
// for the textInfo input parameter.
func createTextInfo(textInfo *m.TextInfo, db *gorm.DB, config cnf.Configurations) (err error) {
	if textInfo.Token != "" {

		_, err = upsertTextInfo(textInfo, db, config)
	} else {
		err = errors.New("no token value found")
	}

	return err
}

// - - - - Batch processing functions - - - -

// BatchCreate addresses the user's request for batch create
func BatchCreate(r *http.Request, db *gorm.DB, config cnf.Configurations) (processStatus []c.TokenProcessStatus, err error) {
	var textInfoList []m.TextInfo

	err = json.NewDecoder(r.Body).Decode(&textInfoList)

	if err != nil {
		return processStatus, err
	}

	for _, textInfo := range textInfoList {

		ps, err := upsertTextInfo(&textInfo, db, config)
		if err != nil {
			ps = c.TokenProcessStatus{Id: textInfo.Id,
				Token:  textInfo.Token,
				Status: c.RequestStatus{Status: "failed", Message: err.Error()}}
		}
		processStatus = append(processStatus, ps)
	}

	return processStatus, nil
}

// BatchUpdate addresses the user's request for batch update
func BatchUpdate(r *http.Request, db *gorm.DB) (processStatus []c.TokenProcessStatus, err error) {
	var textInfoList []m.TextInfo

	err = json.NewDecoder(r.Body).Decode(&textInfoList)

	if err != nil {
		return processStatus, err
	}

	for _, textInfo := range textInfoList {
		ps := updateSingleRecord(&textInfo, db)

		processStatus = append(processStatus, ps)
	}

	return processStatus, nil
}

// BatchDelete addresses the user's request for batch delete
func BatchDelete(r *http.Request, db *gorm.DB) (processStatus []c.TokenProcessStatus, err error) {
	var tiList []m.TextInfoProxy
	err = json.NewDecoder(r.Body).Decode(&tiList)
	if err != nil {
		return processStatus, err
	}

	// 1. Retrieve the ids and create a map from the request input
	tpl := new(m.TextInfoProxyList)
	tpl.SetList(tiList)
	tiRequestMap, requestIds := tpl.List2Map()

	// 2. Find all records matching to the requested Ids in the database
	var storedList []m.TextInfoProxy
	db.Table("text_info").Find(&storedList, requestIds)

	// 3. Retrieve the ids and create a map from the stored records
	tpl.SetList(storedList)
	tiStoredMap, _ := tpl.List2Map()

	// Iterate over the map with the records marked to be deleted and process them one by one
	// Reject read-only records
	// Ignore non-existent records
	// Delete the valid non-readonly records
	for key := range tiRequestMap {

		if storedRecord, ok := tiStoredMap[key]; ok {
			// the id is for an existing record
			if storedRecord.IsReadOnly {
				// do not delete the read-only record
				processStatus = append(processStatus,
					c.TokenProcessStatus{Id: key, Token: storedRecord.Token, Status: c.RequestStatus{Status: "faled", Message: fmt.Sprintf("Readn-only record (id = %d) cannot be deleted", key)}})
			} else {
				db.Exec("delete from text_info where id = ?", key)

				processStatus = append(processStatus,
					c.TokenProcessStatus{Id: key, Token: storedRecord.Token, Status: c.RequestStatus{Status: "deleted"}})
			}
		} else {
			// missing record
			processStatus = append(processStatus,
				c.TokenProcessStatus{Id: key, Token: "", Status: c.RequestStatus{Status: "faled", Message: fmt.Sprintf("No record with id %d", key)}})
		}
	}

	return processStatus, nil
}
