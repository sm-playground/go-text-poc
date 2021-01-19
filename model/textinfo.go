package model

type TextInfo struct {
	// gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Token      string `gorm:"type:varchar(100)"`
	Text       string `gorm:"type:varchar(2000)"`
	Noun       string `json:"noun"`
	Function   string `json:"function"`
	Action     string `json:"action";gorm:"size:30"`
	Module     string `json:"module"`
	Country    string `json:"country"`
	Locale     string `json:"locale"`
	SourceType string `json:"sourceType"`
	SourceId   string `json:"sourceiId"`
	TargetId   string `json:"targetId"`
	ReadOnly   bool   `json:"readOnly"`
}

type TokenText struct {
	Token  string `json:"token"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

type TextInfoPayload struct {
	TargetId       string             `json:"targetId"`
	Locale         string             `json:"locale"`
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
