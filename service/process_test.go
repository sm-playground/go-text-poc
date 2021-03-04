package service

import (
	"github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"strconv"
	"testing"
)

const TEST_TOKEN = "CREATE.RECORD.TYPE.IN.ARMENIAN"

// = = = = = = = = = DeleteTextInfo = = = = = = = =

// TestDeleteTextInfo - targets DeleteTextInfo
func TestDeleteTextInfo(t *testing.T) {
	token := "ABC.DEF.GHI.ABRACADABRA"
	var textInfo = m.TextInfo{
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

// = = = = = = = = = OverwriteTextInfo = = = = = = = =

func TestOverwriteTextInfo(t *testing.T) {
	token := "RECORD.IN.ENGLISH.ABRA.CADABRA"

	tearDown([]string{token})

	var textInfo = m.TextInfo{
		Token:  token,
		Text:   "This is just for testing",
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
	ti.Text = "The text is updated runtime by the test"
	ti.Action = "Update"

	var err error
	ti, err = OverwriteTextInfo(params, ti)
	if err != nil {
		t.Errorf("Failed overwriting the record with id %s", params["id"])
	} else {
		if ti.Action != "Update" || ti.Text != "The text is updated runtime by the test" {
			t.Errorf("Failed overwriting the record %+v", ti)
		}
	}

	// Cleanup at the very end
	tearDown([]string{token})
}

// = = = = = = = = = CreateTextInfo = = = = = = = =

// TestCreateTextInfoFallback verifies the CreateTextInfo - a method called on POST request and
// creates a fallback record
func TestCreateTextInfoFallback(t *testing.T) {

	// delete all test records first
	deleteTestTokens()

	var textInfo = m.TextInfo{
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

	var textInfo = m.TextInfo{
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

	var textInfo = m.TextInfo{
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

// = = = = = = = = = BatchCreate = = = = = = = =

// TestBatchCreate tests creation of multiple records in a batch
func TestBatchCreate(t *testing.T) {
	token := "IA.AR.ARINVOICE.VIEW.LABEL.ABCD.ONE"
	tearDown([]string{token})

	textInfoList := []m.TextInfo{
		{Token: token, Text: "A fallback, Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true},
		{Token: token, Text: "Localized text - Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true, IsFallback: false},
		{Token: token, Text: "Customized text 123456 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "123456", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
		{Token: token, Text: "Customized text 987654 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "987654", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
	}

	response, err := BatchCreate(textInfoList)
	if err != nil {
		t.Errorf("Failed creating a batch of records for a token %s", token)
	} else {
		if response == nil {
			t.Errorf("Incorrect response on batch create")
		} else {
			if len(response) != 4 {
				t.Errorf("Incorrect number of created records on batch create")
			} else {
				for _, r := range response {
					if r.Status.Status != "created" {
						t.Errorf("Failed creating a record on batch create")
					}
				}
			}
		}
	}

	// Clean after yourself
	tearDown([]string{token})
}

// = = = = = = = = = BatchUpdate = = = = = = = =

// TestBatchUpdate
func TestBatchUpdate(t *testing.T) {
	token := "IA.AR.ARINVOICE.VIEW.LABEL.ABCD.ONE"
	tearDown([]string{token})

	textInfoList := []m.TextInfo{
		{Token: token, Text: "A fallback, Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true},
		{Token: token, Text: "Localized text - Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true, IsFallback: false},
		{Token: token, Text: "Customized text 123456 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "123456", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
		{Token: token, Text: "Customized text 987654 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "987654", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
	}

	var idOne, idTwo int

	response, err := BatchCreate(textInfoList)
	if err != nil {
		t.Errorf("Failed creating a batch of records for a token %s", token)
	} else {
		for ind, r := range response {
			if r.Status.Status != "created" {
				t.Errorf("Failed creating a record on batch create")
			} else {
				if ind == 2 {
					idOne = r.Id
				}
				if ind == 3 {
					idTwo = r.Id
				}
			}
		}
	}

	textInfoList = []m.TextInfo{
		{Id: idOne, Token: token, Text: "Updated Customized text 123456 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "123456", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
		{Id: idTwo, Token: token, Text: "Updated Customized text 987654 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "987654", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false, Action: "Update"},
	}
	response, err = BatchUpdate(textInfoList)
	if err != nil {
		t.Errorf("Failed updating a batch of records for a token %s", token)
	} else {
		for _, r := range response {
			if r.Status.Status != "updated" {
				t.Errorf("Failed updating a record on batch update")
			}
		}

	}

	// Clean after yourself
	tearDown([]string{token})
}

// = = = = = = = = = BatchDelete = = = = = = = =

// TestBatchDelete
func TestBatchDelete(t *testing.T) {
	token := "IA.AR.ARINVOICE.VIEW.LABEL.ABCD.ONE"
	tearDown([]string{token})

	textInfoList := []m.TextInfo{
		{Token: token, Text: "A fallback, Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true},
		{Token: token, Text: "Localized text - Գրանցման տեսակը", Locale: "am-AM", IsReadOnly: true, IsFallback: false},
		{Token: token, Text: "Customized text 123456 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "123456", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
		{Token: token, Text: "Customized text 987654 - Գրանցման տեսակը", Locale: "am-AM", TargetId: "987654", SourceId: "PRT-099", IsReadOnly: false, IsFallback: false},
	}

	response, err := BatchCreate(textInfoList)
	if err != nil {
		t.Errorf("Failed creating a batch of records for a token %s", token)
	} else {
		var tiList []m.TextInfoProxy
		for _, ti := range response {
			tiList = append(tiList, m.TextInfoProxy{Id: ti.Id})
		}

		result, e := BatchDelete(tiList)
		if e != nil {
			t.Errorf("Failed properly deleting a batch of records for a token %s", token)
		} else {
			var deletedCounter, rejectedCounter int
			for _, r := range result {
				if r.Status.Status == "deleted" {
					deletedCounter++
				}
				if r.Status.Status == "faled" {
					rejectedCounter++
				}
			}
			if deletedCounter != 2 || rejectedCounter != 2 {
				t.Errorf("Failed properly deleting / rejecting a batch of records for a token %s", token)
			}
		}
	}

	tearDown([]string{token})
}

// = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =

// deleteTestTokens deletes all the records with the token set to test token
//
// Shall be called before every test for cleanup
func deleteTestTokens() {
	connection := db.GetConnection()
	connection.Delete(m.TextInfo{}, "token = ?", TEST_TOKEN)
}
