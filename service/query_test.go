package service

import (
	"github.com/sm-playground/go-text-poc/db"
	m "github.com/sm-playground/go-text-poc/model"
	"strings"
	"testing"
)

const TEST_TARGET_ID = "13579"
const TEST_SOURCE_ID = "PRT-009"

// = = = = = = = = = GetTextInfo = = = = = = = =

// TestGetTextInfo verifies the GetTextInfo - a method called on GET request with no parameters
func TestGetTextInfo(t *testing.T) {
	textInfoList := GetTextInfo(
		nil,
		"")

	if textInfoList == nil || len(textInfoList) == 0 {
		t.Errorf("Failed to query textinfo records")
	}
}

// TestGetTextInfoWithValidToken verifies the GetTextInfo - a method called on GET request with valid token
func TestGetTextInfoWithValidToken(t *testing.T) {
	var tokens = []string{"IA.AR.ARINVOICE.EDIT.LABEL"}
	textInfoList := GetTextInfo(
		tokens,
		"")

	if textInfoList == nil || len(textInfoList) == 0 {
		t.Errorf("Failed to query textinfo records")
	}
}

// TestGetTextInfoWithInvalidToken verifies the GetTextInfo - a method called on GET request with INvalid token
func TestGetTextInfoWithInvalidToken(t *testing.T) {
	var tokens = []string{"I.AM.AN.INVALID.TOKEN"}
	textInfoList := GetTextInfo(
		tokens,
		"")

	if len(textInfoList) > 0 {
		t.Errorf("Failed to query textinfo records")
	}
}

// = = = = = = = = = GetSingleTextInfo = = = = = = = =

// TestGetSingleTextInfo GET with correct input for the record id
func TestGetSingleTextInfo(t *testing.T) {
	params := make(map[string]string)
	params["id"] = "9"

	textInfo, err := GetSingleTextInfo(params)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if textInfo.Id == 0 {
			t.Errorf("The record with %s id cannot be found", params["id"])
		}
	}
}

// TestGetSingleTextInfo GET with INcorrect input for the record id
func TestGetSingleTextInfoInvalidID(t *testing.T) {
	params := make(map[string]string)
	params["id"] = "999"

	textInfo, err := GetSingleTextInfo(params)
	if err == nil {
		t.Errorf(err.Error())
	}
	if textInfo.Id != 0 {
		t.Errorf("Found the record with id %s", params["id"])
	}
}

// = = = = = = = = = ReadData = = = = = = = =

// TestReadDataSingleToken test reading a single token
func TestReadDataSingleToken(t *testing.T) {
	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
	)

	err := setup()
	if err != nil {
		t.Errorf(err.Error())
	}

	var payload m.TextInfoPayload = m.TextInfoPayload{
		SourceId: TEST_SOURCE_ID,
		TargetId: TEST_TARGET_ID,
		Locale:   "en-US",
		Tokens: []m.TokenPayload{
			{Token: "RECORD.TYPE.IN.ENGLISH"},
		},
	}

	data, err := ReadData(payload)
	if err != nil {
		t.Errorf("Failed reading data")
	}

	if len(data) != 1 {
		t.Errorf("Expected a single record, got %d records", len(data))
	}

	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
	)
}

// TestReadDataMultipleTokens test reading multiple tokens
func TestReadDataMultipleTokens(t *testing.T) {
	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
	)

	err := setup()
	if err != nil {
		t.Errorf(err.Error())
	}

	var payload m.TextInfoPayload = m.TextInfoPayload{
		SourceId: TEST_SOURCE_ID,
		TargetId: TEST_TARGET_ID,
		Locale:   "en-US",
		Tokens: []m.TokenPayload{
			{Token: "RECORD.TYPE.IN.ENGLISH"},
			{Token: "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
		},
	}

	data, err := ReadData(payload)
	if err != nil {
		t.Errorf("Failed reading data")
	}

	if len(data) != 2 {
		t.Errorf("Expected two records, got %d records", len(data))
	}

	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
	)
}

func TestReadDataParametrizedToken(t *testing.T) {
	token := "PARAMETRIZED.RECORD.IN.ENGLISH"
	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH",
		token},
	)

	err := setup()
	if err != nil {
		t.Errorf(err.Error())
	}

	var textInfo m.TextInfo = m.TextInfo{
		Token:  token,
		Text:   "The page {PAGE} cannot be {ACTION}",
		Locale: "en-US",
		Action: "Test",
	}

	if createTokenRecords(textInfo, "en-US") != nil {
		t.Errorf("Failed creating records for tokenb - %s", token)
	}

	var payload m.TextInfoPayload = m.TextInfoPayload{
		SourceId: TEST_SOURCE_ID,
		TargetId: TEST_TARGET_ID,
		Locale:   "en-US",
		Tokens: []m.TokenPayload{
			{
				Token: token,
				Placeholders: []m.TokenPlaceholder{
					{Name: "PAGE", Value: "Introduction"},
					{Name: "ACTION", Value: "Deleted"},
				},
			},
			{Token: "ONE.MORE.RECORD.TYPE.IN.ENGLISH"},
		},
	}

	data, err := ReadData(payload)
	if err != nil {
		t.Errorf("Failed reading data")
	}

	if len(data) != 2 {
		t.Errorf("Expected two records, got %d records", len(data))
	}

	if data[0].Token == token {
		if !strings.Contains(data[0].Text, "The page Introduction cannot be Deleted") {
			t.Errorf("Parametrized token %s does not contain expected string %s", data[0].Text, "The page Introduction cannot be Deleted")
		}
	}

	tearDown([]string{"RECORD.TYPE.IN.ARMENIAN", "RECORD.TYPE.IN.ENGLISH",
		"ANOTHER.RECORD.TYPE.IN.ENGLISH", "ONE.MORE.RECORD.TYPE.IN.ENGLISH",
		token},
	)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

// tearDown - tear down routine
func tearDown(tokens []string) {
	for _, token := range tokens {
		deleteTokenRecords(token)
	}
}

// setup - setup routine
func setup() error {
	var textInfo m.TextInfo = m.TextInfo{Token: "RECORD.TYPE.IN.ARMENIAN", Action: "Test"}
	err := createTokenRecords(textInfo, "am-AM")
	if err != nil {
		return err
	}

	textInfo = m.TextInfo{Token: "RECORD.TYPE.IN.ENGLISH", Action: "Test"}
	err = createTokenRecords(textInfo, "en-US")
	if err != nil {
		return err
	}

	textInfo = m.TextInfo{Token: "ONE.MORE.RECORD.TYPE.IN.ENGLISH", Action: "Test"}
	err = createTokenRecords(textInfo, "en-US")
	if err != nil {
		return err
	}

	textInfo = m.TextInfo{Token: "ANOTHER.RECORD.TYPE.IN.ENGLISH", Action: "Test"}
	err = createTokenRecords(textInfo, "en-US")
	if err != nil {
		return err
	}

	return nil
}

// createTokenRecords A routine that creates three records for a given token: fallback, localized, and customized
func createTokenRecords(textInfo m.TextInfo, locale string) (err error) {

	txt := ""
	if textInfo.Text != "" {
		txt = textInfo.Text + " "
	}

	textInfo.IsFallback = true
	textInfo.Locale = locale

	// create fallback
	textInfo.Text = txt + locale + " " + textInfo.Token + " Fallback"
	_, err = CreateTextInfo(textInfo)
	if err != nil {
		return err
	}

	textInfo.IsFallback = false
	textInfo.Text = txt + locale + " " + textInfo.Token + " Localized"
	_, err = CreateTextInfo(textInfo)
	if err != nil {
		return err
	}

	textInfo.IsFallback = false
	textInfo.TargetId = TEST_TARGET_ID
	textInfo.SourceId = TEST_SOURCE_ID
	textInfo.Text = txt + locale + " " + textInfo.Token + " " + textInfo.TargetId + " " + textInfo.SourceId + " Customized"
	_, err = CreateTextInfo(textInfo)
	if err != nil {
		return err
	}

	return nil
}

// deleteTokenRecords delete all existing records for corresponding to a specified token
func deleteTokenRecords(token string) {
	connection := db.GetConnection()
	connection.Delete(m.TextInfo{}, "token = ?", token)
}
