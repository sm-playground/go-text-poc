package service

import (
	"encoding/json"
	"errors"
	"fmt"
	cache "github.com/sm-playground/go-text-poc/cache_client"
	c "github.com/sm-playground/go-text-poc/common"
	cnf "github.com/sm-playground/go-text-poc/config"
	d "github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"io"
	"strconv"
)

// DeleteTextInfo Deletes a single record from the text_info table
func DeleteTextInfo(params map[string]string) (textInfo m.TextInfo, err error) {

	db := d.GetConnection()

	var requestStatus c.RequestStatus

	var deletedRecordId = params["id"]
	if db.First(&textInfo, deletedRecordId).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", deletedRecordId))

	} else {
		err = nil
		db.Delete(&textInfo)

		if cnf.GetInstance().Get().Cache.UseCache {
			cacheClient, _ := cache.GetCacheClient()
			_ = cacheClient.Invalidate("*" + textInfo.Token + "*")
		}

		requestStatus.Status = "success"
		requestStatus.Message = fmt.Sprintf("The record with id=%s was deleted", deletedRecordId)
	}

	return textInfo, err
}

// CreateTextInfo Creates a single record in the text_info table
func CreateTextInfo(textInfo m.TextInfo) (m.TextInfo, error) {
	config := cnf.GetInstance().Get()

	if textInfo.TargetId == "" {
		// No target Id is specified so this record is applicable to all targets
		// Hence, the source shall be set to the service owner source Id
		textInfo.SourceId = config.ServiceOwnerSourceId
	}

	err := createTextInfo(&textInfo)
	if err == nil {
		queryInput := SingleQueryInput{
			ServiceOwnerSourceId: config.ServiceOwnerSourceId,
			DefaultCountry:       config.DefaultCountry,
			DefaultLanguage:      config.DefaultLanguage,
			Country:              textInfo.Country,
			Language:             textInfo.Language,
			TargetId:             textInfo.TargetId,
			SourceId:             textInfo.SourceId}

		useCache := cnf.GetInstance().Get().Cache.UseCache
		if useCache {
			cacheKey := getReadDataCacheKey(queryInput, textInfo.Token+"%")
			var cacheClient cache.CacheClient
			if cacheClient, err = cache.GetCacheClient(); err == nil {
				_ = cacheClient.Set(cacheKey, textInfo)
			}
		}
	}

	return textInfo, err
}

// UpdateTextInfo updates the record in the text_info table. Only the fields specified in the request are updated
func UpdateTextInfo(params map[string]string, body io.ReadCloser) (textInfo m.TextInfo, err error) {

	db := d.GetConnection()

	var id = params["id"]
	if db.First(&textInfo, id).RecordNotFound() {
		err = errors.New(fmt.Sprintf("The record with id=%s is not found", id))
		return textInfo, err
	} else {
		var newPyload m.TextInfo

		err = json.NewDecoder(body).Decode(&newPyload)
		if err != nil {
			return textInfo, err
		}

		textInfo.Overwrite(newPyload)

		db.Save(&textInfo)

		return textInfo, nil
	}
}

// OverwriteTextInfo updates the record in the text_info table. ALL fields are updated
func OverwriteTextInfo(params map[string]string, ti m.TextInfo) (m.TextInfo, error) {

	db := d.GetConnection()

	var textInfo m.TextInfo
	var id = params["id"]
	if db.First(&textInfo, id).RecordNotFound() {
		err := errors.New(fmt.Sprintf("The record with id=%s is not found", id))
		return textInfo, err
	} else {
		if v, err := strconv.Atoi(params["id"]); err == nil {
			// Convert the value of Id into integer
			ti.Id = v
			db.Save(&ti)
		}

		return ti, nil
	}
}

// createTextInfo creates a single record in the text_info table
// for the textInfo input parameter.
func createTextInfo(textInfo *m.TextInfo) (err error) {
	if textInfo.Token != "" {

		_, err = upsertTextInfo(textInfo)
	} else {
		err = errors.New("no token value found")
	}

	return err
}

// - - - - Batch processing functions - - - -

// BatchCreate addresses the user's request for batch create
func BatchCreate(textInfoList []m.TextInfo) (processStatus []c.TokenProcessStatus, err error) {

	for _, textInfo := range textInfoList {

		ps, err := upsertTextInfo(&textInfo)
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
func BatchUpdate(textInfoList []m.TextInfo) (processStatus []c.TokenProcessStatus, err error) {

	for _, textInfo := range textInfoList {
		ps := updateSingleRecord(&textInfo)

		processStatus = append(processStatus, ps)
	}

	return processStatus, nil
}

// BatchDelete addresses the user's request for batch delete
func BatchDelete(tiList []m.TextInfoProxy) (processStatus []c.TokenProcessStatus, err error) {
	db := d.GetConnection()

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
