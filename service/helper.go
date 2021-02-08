package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	c "github.com/sm-playground/go-text-poc/common"
	cnf "github.com/sm-playground/go-text-poc/config"
	m "github.com/sm-playground/go-text-poc/model"
)

// upsertTextInfo
func upsertTextInfo(textInfo *m.TextInfo, db *gorm.DB, config cnf.Configurations) (processStatus c.TokenProcessStatus, err error) {
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
