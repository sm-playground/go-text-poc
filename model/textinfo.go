package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	cache "github.com/sm-playground/go-text-poc/cache_client"
	"reflect"
)

type TextInfo struct {
	// gorm.Model
	Id         int    `gorm:"autoIncrement:true" json:"id"`
	Token      string `gorm:"type:varchar(100)" json:"token"`
	Text       string `gorm:"type:varchar(2000)" json:"text"`
	TargetId   string `json:"targetId"`
	SourceId   string `json:"sourceId"`
	Language   string `json:"language"`
	Country    string `json:"country"`
	Locale     string `gorm:"-" json:"-"`
	Action     string `gorm:"size:30" json:"action"`
	SourceType string `json:"sourceType"`
	IsReadOnly bool   `gorm:"column:read_only" json:"readOnly"`
	IsFallback bool   `gorm:"column:is_fallback" json:"fallBack"`

	Temp string `gorm:"<-:false" json:"-"`

	Noun        string `json:"noun"`
	Function    string `json:"function"`
	Verb        string `json:"verb"`
	Application string `json:"application"`
	Module      string `json:"module"`
}

// Unmarshal convert the passed parameter into the TextInfo object
func (ti *TextInfo) Unmarshal(data interface{}) (textInfo TextInfo, err error) {
	if data == nil {
		return textInfo, errors.New(cache.NIL_VALUE_ERROR_MESSAGE)
	}

	js, _ := redis.String(data, nil)

	err = json.Unmarshal([]byte(js), &textInfo)

	return textInfo, err

}

// Overwrite overwrites all fields of the object with the values from the method input object
func (ti *TextInfo) Overwrite(textInfo TextInfo) {
	// myValue := reflect.ValueOf(ti)
	newObjectValue := reflect.ValueOf(textInfo)
	numberOfFields := newObjectValue.NumField()
	for i := 0; i < numberOfFields; i++ {
		fmt.Printf("%d.Type:%T || Value:%#v\n",
			i+1, newObjectValue.Field(i), newObjectValue.Field(i))

		fmt.Println("Kind is ", newObjectValue.Field(i).Kind())
	}
}

type TextInfoPayload struct {
	TargetId       string             `json:"targetId"`
	SourceId       string             `json:"sourceId"`
	Locale         string             `json:"locale"`
	Language       string             `json:"language"`
	Country        string             `json:"country"`
	ResponseFormat TextResponseFormat `json:"format"`
	Tokens         []TokenPayload     `json:"tokens"`
}

type TextResponseFormat struct {
	DateFormat     string `json:"date"`
	TimeFormat     string `json:"time"`
	NumberFormat   string `json:"number"`
	CurrencyFormat string `json:"currency"`
	CurrencySymbol string `json:"currencySymbol"`
}

type TokenPayload struct {
	Token        string             `json:"token"`
	Placeholders []TokenPlaceholder `json:"placeholders"`
}

type TokenPlaceholder struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Obfuscate bool   `json:"obfuscate"`
}

// - - - - - - - - - - - TextInfoList - - - - - - - - - - -
type TextInfoList struct {
	List []TextInfo
}

// SetList sets the collection to the wrapper structure
func (proxyList *TextInfoList) SetList(list []TextInfo) {
	proxyList.List = list
}

// Append adds the single element of the TextInfo object to the list
func (proxyList *TextInfoList) Append(ti TextInfo) {
	proxyList.List = append(proxyList.List, ti)
}

// Unmarshal convert the passed parameter into the TextInfoList
func (proxyList *TextInfoList) Unmarshal(data interface{}) (til TextInfoList, err error) {

	if data == nil {
		return til, errors.New(cache.NIL_VALUE_ERROR_MESSAGE)
	}

	js, _ := redis.String(data, nil)

	err = json.Unmarshal([]byte(js), &til)

	return til, err
}

// - - - - - - - - - - - TextInfoProxyList - - - - - - - - - - -
// a structure representing the TextInfoProxy list
type TextInfoProxyList struct {
	List []TextInfoProxy
}

// SetList sets the collection to the wrapper structure
func (proxyList *TextInfoProxyList) SetList(list []TextInfoProxy) {
	proxyList.List = list
}

// List2Map converts the structure representing the list of the TextInfoProxy objects
// into the map where the key of the item in the map is the Id of the TextInfoProxy
func (proxyList *TextInfoProxyList) List2Map() (map[int]TextInfoProxy, []int) {

	var ids []int
	m := make(map[int]TextInfoProxy)

	for _, ti := range proxyList.List {
		m[ti.Id] = ti
		ids = append(ids, ti.Id)
	}

	return m, ids
}

// - - - - - - - - - - - TokenText - - - - - - - - - - -
type TokenText struct {
	Token  string `json:"token"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

type TokenTextList struct {
	List []TokenText
}

// SetList sets the collection to the wrapper structure
func (proxyList *TokenTextList) SetList(list []TokenText) {
	proxyList.List = list
}

// Append adds the single element of the TextInfo object to the list
func (proxyList *TokenTextList) Append(ti TokenText) {
	proxyList.List = append(proxyList.List, ti)
}

// Unmarshal convert the passed parameter into the TextInfoList
func (proxyList *TokenTextList) Unmarshal(data interface{}) (ttl TokenTextList, err error) {

	if data == nil {
		return ttl, errors.New(cache.NIL_VALUE_ERROR_MESSAGE)
	}

	js, _ := redis.String(data, nil)

	err = json.Unmarshal([]byte(js), &ttl)

	return ttl, err
}

// - - - - - - - - - - - TextInfoProxy - - - - - - - - - - -
// The proxy object containing limited fields from the Textinfo
type TextInfoProxy struct {
	Id         int    `json:"id"`
	Token      string `json:"token"`
	IsReadOnly bool   `gorm:"column:read_only" json:"readOnly"`
}
