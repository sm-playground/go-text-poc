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
	/*
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
				return errors.New("the fallback record is already in the database and tries to insert another one")
			}
		}
	*/

	if textInfo.Token != "" {

		// db.Create(&textInfo)
		_, err = upsertTextInfo(textInfo, db, config)
	} else {
		err = errors.New("no token value found")
	}

	return err
}

// - - - - Batch processing functions - - - -
// BatchCreate addresses the user's request for batch update
func BatchCreate(r *http.Request, db *gorm.DB, config cnf.Configurations) (processStatus []c.TokenProcessStatus, err error) {
	var textInfoList []m.TextInfo

	err = json.NewDecoder(r.Body).Decode(&textInfoList)

	if err != nil {
		return processStatus, err
	}

	// createTextInfoList, updateTextInfoList := preprocessForCreate(textInfoList, config, db)

	for _, textInfo := range textInfoList {

		ps, _ := upsertTextInfo(&textInfo, db, config)
		processStatus = append(processStatus, ps)

		/*
			if textInfo.Language == "" {
				textInfo.Language = config.DefaultLanguage
			}
			if textInfo.Country == "" {
				textInfo.Language = config.DefaultCountry
			}

			var count int64
			db.Model(&m.TextInfo{}).
				Where("token = ? AND language = ? AND country is NULL AND is_fallback",
					textInfo.Token, textInfo.Language, textInfo.Country).
				Count(&count)
			if count == 0 {
				// No fallback is found for the token / locale.
				// Make it a fallback record and set the service owner source Id
				textInfo.IsFallback = true
				textInfo.SourceId = config.ServiceOwnerSourceId
				textInfo.TargetId = ""
				db.Create(&textInfo)
				processStatus = append(processStatus, c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}})
			} else {
				// Fallback is found. Check for the localized record
				db.Model(&m.TextInfo{}).
					Where("token = ? AND language = ? AND country = ? AND source_id = ? AND NOT is_fallback",
						textInfo.Token, textInfo.Language, textInfo.Country, config.ServiceOwnerSourceId).
					Count(&count)
				if count == 0 {
					// no localized version
					textInfo.IsFallback = false
					textInfo.SourceId = config.ServiceOwnerSourceId
					textInfo.TargetId = ""
					db.Create(&textInfo)
					processStatus = append(processStatus, c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}})
				} else {
					// neither fallback nor localized. Check for customized record
					var tiList []m.TextInfo
					db.Where("token = ? AND language = ? AND country = ? AND source_id != ? AND target_id = ? AND NOT is_fallback",
						textInfo.Token, textInfo.Language, textInfo.Country, config.ServiceOwnerSourceId, textInfo.TargetId).Find(&tiList)
					if len(tiList) > 1 {
						// TODO ERROR condition.
					} else if len(tiList) == 0 {
						// mark as create
						db.Create(&textInfo)
						processStatus = append(processStatus, c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}})
					} else {
						// mark as update
						textInfo.Id = tiList[0].Id
						db.Save(&textInfo)
						processStatus = append(processStatus, c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "updated"}})
					}
				}
			}
		*/

	}

	return processStatus, nil
}
