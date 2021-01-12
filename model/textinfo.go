package model

type TextInfo struct {
	// gorm.Model
	Id         int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Token      string `gorm:"type:varchar(100)"`
	Text       string `gorm:"type:varchar(2000)"`
	Noun       string
	Function   string
	Action     string `gorm:"size:30"`
	Module     string
	Country    string
	Locale     string
	SourceType string
	SourceId   string
	TargerId   string
	ReadOnly   bool
}
