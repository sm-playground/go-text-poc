package model

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

type TokenText struct {
	Token  string `json:"token"`
	Text   string `json:"text"`
	Status string `json:"status"`
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
