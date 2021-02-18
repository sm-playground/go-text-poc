package service

import (
	"errors"
	"fmt"
	c "github.com/sm-playground/go-text-poc/common"
	cnf "github.com/sm-playground/go-text-poc/config"
	d "github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"strings"
)

// upsertTextInfo creates the record or updates it (overwrites) if the matching one is identified
//
// 1. Checks for the fallback record and create it if it is not found
//
// 2. Checks for the localized record and create it if it is not found
//
// 3. Tries to identifies the newly supplied data with the existing record in the database. If found, updates it.
// Otherwise, creates the new customized record
func upsertTextInfo(textInfo *m.TextInfo) (processStatus c.TokenProcessStatus, err error) {
	config := cnf.GetInstance().Get()

	db := d.GetConnection()

	if textInfo.Locale != "" {
		language, country, err := retrieveLocale(textInfo.Locale)
		if err != nil {
			errorMessage := fmt.Sprintf("Locale is passed in incorrect format: %s", textInfo.Locale)
			processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "failed", Message: errorMessage}}
			err = errors.New(errorMessage)
			return processStatus, err
		} else {
			textInfo.Language = language
			textInfo.Country = country
		}
	}

	if textInfo.Language == "" {
		textInfo.Language = config.DefaultLanguage
	}
	if textInfo.Country == "" {
		textInfo.Language = config.DefaultCountry
	}

	// Check if there is a fallback record for the specified language
	var count int64
	db.Model(&m.TextInfo{}).
		Where("token = ? AND language = ? AND country = '' AND is_fallback",
			textInfo.Token, textInfo.Language).
		Count(&count)
	if count == 0 {
		// No fallback is found for the token / locale.
		// Make it a fallback record and set the service owner source Id
		textInfo.IsFallback = true
		textInfo.Country = ""
		textInfo.SourceId = config.ServiceOwnerSourceId
		textInfo.TargetId = ""
		db.Create(&textInfo)
		processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}}
	} else {
		// Fallback is found. Check for the localized record
		db.Model(&m.TextInfo{}).
			Where("token = ? AND language = ? AND country = ? AND source_id = ? AND NOT is_fallback",
				textInfo.Token, textInfo.Language, textInfo.Country, config.ServiceOwnerSourceId).
			Count(&count)
		if count == 0 {
			// no localized version
			if textInfo.TargetId != "" {
				// The user is trying to create a customized record with the specific target Id
				// while there is no localized version found for that locale
				errorMessage := fmt.Sprintf("Cannot create a customized record with the specific target (%s) without localized version - %s-%s",
					textInfo.TargetId, textInfo.Language, textInfo.Country)
				processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "failed", Message: errorMessage}}
				err = errors.New(errorMessage)
			} else {
				textInfo.IsFallback = false
				textInfo.SourceId = config.ServiceOwnerSourceId
				db.Create(&textInfo)
				processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}}
			}

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
				processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "created"}}
			} else {
				// mark as update
				textInfo.Id = tiList[0].Id
				db.Save(&textInfo)
				processStatus = c.TokenProcessStatus{Id: textInfo.Id, Token: textInfo.Token, Status: c.RequestStatus{Status: "updated"}}
			}
		}
	}

	return processStatus, err
}

// updateSingleRecord validates the passed TextInfo object and updates the matching
// record in the database.
func updateSingleRecord(textInfo *m.TextInfo) c.TokenProcessStatus {
	db := d.GetConnection()

	var ps c.TokenProcessStatus
	if textInfo.Id == 0 {
		// Error condition. Id is required for update
		ps = c.TokenProcessStatus{Id: textInfo.Id,
			Token:  textInfo.Token,
			Status: c.RequestStatus{Status: "failed", Message: "Id is required for update"}}
	} else {
		var ti m.TextInfo
		if db.First(&ti, textInfo.Id).RecordNotFound() {
			// Find a record in the database with the given Id
			ps = c.TokenProcessStatus{Id: textInfo.Id,
				Token:  textInfo.Token,
				Status: c.RequestStatus{Status: "failed", Message: fmt.Sprintf("No record with Id -> %d", textInfo.Id)}}
		} else {
			if ti.IsReadOnly {
				ps = c.TokenProcessStatus{Id: textInfo.Id,
					Token:  textInfo.Token,
					Status: c.RequestStatus{Status: "failed", Message: "Cannot update read-only field"}}
			} else {
				db.Save(&textInfo)
				ps = c.TokenProcessStatus{Id: textInfo.Id,
					Token:  textInfo.Token,
					Status: c.RequestStatus{Status: "updated"}}
			}
		}
	}

	return ps
}

// retrieveLocale retrieved the language and country from the locale passed in the format of en-US
func retrieveLocale(locale string) (language string, country string, err error) {
	l := strings.Split(locale, "-")

	if len(l) != 2 {
		// the locale is expected in the format of "en-US"
		// Ignore all other cases
		return language, country, errors.New(fmt.Sprintf("incorrect format for locale: %s", locale))
	}
	language = strings.TrimSpace(l[0])
	country = strings.TrimSpace(l[1])

	return language, country, nil

}
