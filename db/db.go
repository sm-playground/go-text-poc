package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	c "github.com/sm-playground/go-text-poc/config"
	m "github.com/sm-playground/go-text-poc/model"
)

var (
	textInfo = []m.TextInfo{
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Record type", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Record type", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.ENTITY", Text: "Customer entity", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.ENTITY", Text: "Customer entity", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "View", Country: "US", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "Edit", Country: "US", Locale: "en", ReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Record type", Action: "View", Country: "CA", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Type d'enregistrement", Action: "View", Country: "CA", Locale: "fr", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Customer name", Action: "View", Country: "CA", Locale: "en", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Nom du client", Action: "View", Country: "CA", Locale: "fr", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "ომხმარებლის სახელი", Action: "View", Country: "GE", Locale: "ge", ReadOnly: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.CUSTOMERNAME", Text: "Հաճախորդի անունը", Action: "View", Country: "AM", Locale: "am", ReadOnly: true},
		// {Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Գրանցման տեսակը", Action: "View", Country: "AM", Locale: "am", ReadOnly: true},
	}
)

// initDatabase initializes the database
// - connect to postgres
// - runs AutoMigrate to create the database
// - populate the text_info table with the localized records list
func InitDatabase(config c.Configurations) (db *gorm.DB) {

	var err error

	dbConnect := fmt.Sprintf("port=%d user=%s dbname=%s sslmode=%s password=%s",
		config.Database.Port,
		config.Database.DBUser,
		config.Database.DBName,
		config.Database.SSLMode,
		config.Database.DBPassword)

	db, err = gorm.Open(config.Database.Dialect, dbConnect)
	if err != nil {
		panic("failed to connect database")
	}

	db.DB().SetMaxIdleConns(config.Database.DBCP.MaxIdle)
	db.DB().SetMaxOpenConns(config.Database.DBCP.MaxActive)

	// set the log mode to see the queries executed by the gorm
	db.LogMode(true)

	// instructs gorm to create the tables with the singular names
	db.SingularTable(true)

	// var initDatabase

	if config.InitDatabase {
		// create tables based on specified structs
		db.AutoMigrate(&m.TextInfo{})

		for index := range textInfo {
			db.Create(&textInfo[index])
		}
	}

	return db
}
