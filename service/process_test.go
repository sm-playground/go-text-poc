package service

import (
	"github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"strconv"
	"testing"
)

const TEST_TOKEN = "CREATE.RECORD.TYPE.IN.ARMENIAN"

// = = = = = = = = = CreateTextInfo = = = = = = = =

// TestCreateTextInfoFallback verifies the CreateTextInfo - a method called on POST request and
// creates a fallback record
func TestCreateTextInfoFallback(t *testing.T) {

	// delete all test records first
	deleteTestTokens()

	var textInfo m.TextInfo = m.TextInfo{
		Token:      TEST_TOKEN,
		Text:       "Գրանցման տեսակը",
		Action:     "Create",
		Locale:     "am-AM",
		IsReadOnly: true,
	}

	ti, err := CreateTextInfo(textInfo)

	if err != nil {
		t.Errorf("Failed creating fallback record %+v", textInfo)
	} else {
		if !ti.IsFallback {
			t.Errorf("Failed creating fallback record %+v", ti)
		}
	}
}

// TestCreateTextInfoLocalized verifies the CreateTextInfo - a method called on POST request and
// creates a localized record
func TestCreateTextInfoLocalized(t *testing.T) {
	// delete all test records first
	deleteTestTokens()

	var textInfo m.TextInfo = m.TextInfo{
		Token:      TEST_TOKEN,
		Text:       "Գրանցման տեսակը",
		Action:     "Create",
		Locale:     "am-AM",
		IsReadOnly: true,
	}

	// Create a fallback record
	ti, err := CreateTextInfo(textInfo)
	if err != nil {
		t.Errorf("Failed creating fallback record %+v", textInfo)
	}

	textInfo.IsFallback = false
	ti, err = CreateTextInfo(textInfo)

	if err != nil {
		t.Errorf("Failed creating localized record %+v", textInfo)
	} else {
		if ti.IsFallback {
			t.Errorf("Failed creating localized record %+v", ti)
		}
	}
}

// TestCreateTextInfoCusomized verifies the CreateTextInfo - a method called on POST request and
// creates a customized record
func TestCreateTextInfoCustomized(t *testing.T) {
	// delete all test records first
	deleteTestTokens()

	var textInfo m.TextInfo = m.TextInfo{
		Token:      TEST_TOKEN,
		Text:       "Գրանցման տեսակը",
		Action:     "Create",
		Locale:     "am-AM",
		IsReadOnly: true,
	}

	// Create a fallback record
	ti, err := CreateTextInfo(textInfo)
	if err != nil {
		t.Errorf("Failed creating fallback record %+v", textInfo)
	}

	textInfo.IsFallback = false
	ti, err = CreateTextInfo(textInfo)

	if err != nil {
		t.Errorf("Failed creating localized record %+v", textInfo)
	} else {
		if ti.IsFallback {
			t.Errorf("Failed creating localized record %+v", ti)
		}
	}

	textInfo.IsFallback = false
	textInfo.TargetId = "13579"
	textInfo.SourceId = "PRT-009"
	ti, err = CreateTextInfo(textInfo)
	if err != nil {
		t.Errorf("Failed creating customized record %+v", textInfo)
	} else {
		if ti.TargetId != "13579" && ti.SourceId != "PRT-009" {
			t.Errorf("Failed creating customized record %+v", textInfo)
		}
	}

	// Clean after test
	deleteTestTokens()

}

// TestDeleteTextInfo - targets DeleteTextInfo
func TestDeleteTextInfo(t *testing.T) {
	token := "ABC.DEF.GHI.ABRACADABRA"
	var textInfo m.TextInfo = m.TextInfo{
		Token:  token,
		Text:   "Abracadabra",
		Locale: "en-US",
		Action: "Test",
	}

	if createTokenRecords(textInfo, "en-US") != nil {
		t.Errorf("Failed creating records for tokenb - %s", token)
	}

	var ti m.TextInfo
	db.GetConnection().Last(&ti)

	params := make(map[string]string)
	params["id"] = strconv.Itoa(ti.Id)

	textInfo, err := DeleteTextInfo(params)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if textInfo.Id == 0 {
			t.Errorf("The record with %s id cannot be found", params["id"])
		}
	}

	deleteTokenRecords(token)
}

func TestDeleteTextInfoInvalid(t *testing.T) {
	params := make(map[string]string)
	params["id"] = "999"

	_, err := DeleteTextInfo(params)
	if err == nil {
		t.Errorf("The test failed, non-existent record with 999 Id was deleted")
	}
}

// = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =

// deleteTestTokens deletes all the records with the token set to test token
//
// Shall be called before every test for cleanup
func deleteTestTokens() {
	connection := db.GetConnection()
	connection.Delete(m.TextInfo{}, "token = ?", TEST_TOKEN)
}
