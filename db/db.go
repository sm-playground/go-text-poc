package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	c "github.com/sm-playground/go-text-poc/config"
	m "github.com/sm-playground/go-text-poc/model"
)

var (
	textInfo = []m.TextInfo{
		// fallback records
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Due in", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Dû en jours", TargetId: "", SourceId: "", Language: "fr", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},

		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Record type", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "View", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Type d'enregistrement", TargetId: "", SourceId: "", Language: "fr", Country: "", Action: "View", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Tipo de registro", TargetId: "", SourceId: "", Language: "es", Country: "", Action: "View", IsReadOnly: true, IsFallback: true},

		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Record type", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Edit Type d'enregistrement", TargetId: "", SourceId: "", Language: "fr", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "EDIT Tipo de registro", TargetId: "", SourceId: "", Language: "es", Country: "", Action: "View", IsReadOnly: true, IsFallback: true},

		// fallback - English only
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.ENTITY", Text: "Customer entity", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "Customer name", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DELIVERY_OPTIONS", Text: "Customer delivery options", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DOC_NUMBER", Text: "Document number", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},
		// fallback parametrized
		{Token: "IA.AR.ARINVOICE.EDIT.HELP.RECORDTYPE", Text: "Help message for {RECORD_TYPE}", TargetId: "", SourceId: "", Language: "en", Country: "", Action: "Edit", IsReadOnly: true, IsFallback: true},

		// localized records
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Localized Due in days", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Localized Dû en jours", TargetId: "", SourceId: "", Language: "fr", Country: "FR", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Localized record type", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "View", IsReadOnly: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Le type d'enregistrement", TargetId: "", SourceId: "", Language: "fr", Country: "FR", Action: "View", IsReadOnly: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Localized record type", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Localized Type d'enregistrement", TargetId: "", SourceId: "", Language: "fr", Country: "FR", Action: "Edit", IsReadOnly: true, IsFallback: false},
		// localized records - English only
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.ENTITY", Text: "The customer entity", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "The customer name", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		// localized parametrized
		{Token: "IA.AR.ARINVOICE.EDIT.HELP.RECORDTYPE", Text: "Localized message for {RECORD_TYPE}", TargetId: "", SourceId: "", Language: "en", Country: "US", Action: "Edit", IsReadOnly: true, IsFallback: false},

		// Customized records
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Customized Type d'enregistrement", TargetId: "123456", SourceId: "PRT-002", Language: "fr", Country: "FR", Action: "Edit", IsReadOnly: true, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.RECORDTYPE", Text: "Customized record type", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.ENTITY", Text: "Customized customer entity", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "Customized customer name", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Customized When due", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DUE_IN_DAYS", Text: "Customized Dû en jours", TargetId: "123456", SourceId: "PRT-001", Language: "fr", Country: "FR", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DELIVERY_OPTIONS", Text: "Customer delivery options", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.DOC_NUMBER", Text: "The document number", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: false, IsFallback: false},
		// customized parametrized
		{Token: "IA.AR.ARINVOICE.EDIT.HELP.RECORDTYPE", Text: "Customized message for {RECORD_TYPE}", TargetId: "123456", SourceId: "PRT-001", Language: "en", Country: "US", Action: "Edit", IsReadOnly: true, IsFallback: false},

		// Customized records without localized records for country but with fallback for language
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "en_CA Customer name", TargetId: "123456", SourceId: "PRT-001", Action: "Edit", Country: "CA", Language: "en", IsReadOnly: false, IsFallback: false},
		// country -> CA, no localized
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "en_CA View record type", TargetId: "123456", SourceId: "PRT-001", Action: "View", Country: "CA", Language: "en", IsReadOnly: true, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "fr_CA Type d'enregistrement", TargetId: "123456", SourceId: "PRT-001", Action: "View", Country: "CA", Language: "fr", IsReadOnly: true, IsFallback: false},

		// Customized records without localized AND fallbackrecords
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "fr_CA CNom du client", TargetId: "123456", SourceId: "PRT-001", Action: "Edit", Country: "CA", Language: "fr", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "ge_GE ომხმარებლის სახელი", TargetId: "123456", SourceId: "PRT-001", Action: "Edit", Country: "GE", Language: "ge", IsReadOnly: false, IsFallback: false},
		{Token: "IA.AR.ARINVOICE.EDIT.LABEL.CUSTOMERNAME", Text: "am_AM Հաճախորդի անունը", TargetId: "123456", SourceId: "PRT-001", Action: "Edit", Country: "AM", Language: "am", IsReadOnly: false, IsFallback: false},

		//
		{Token: "US.PLACEHOLDER.COUNT.FOUR", Text: "The {TOTAL} of {PAYMENT} for {PRODUCT} is {AMOUNT}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "123456", IsFallback: false},
		{Token: "US.PLACEHOLDER.COUNT.ONE", Text: "The life is just {LIFE_DESCRIPTION}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "123456", IsFallback: false},
		{Token: "US.PLACEHOLDER.COUNT.TWO", Text: "Mr. {LEADER}, tear down this {STRUCTURE}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "123456", IsFallback: false},

		{Token: "CONSTRUCTION.ACTION.STRUCTURE", Text: "Dear Mr. {LEADER}, {ACTION} this {STRUCTURE}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "", IsFallback: true},
		{Token: "CONSTRUCTION.ACTION.STRUCTURE", Text: "Mr. {LEADER}, {ACTION} this {STRUCTURE}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "123456", IsFallback: false},
		{Token: "CONSTRUCTION.ACTION.STRUCTURE", Text: "Hey {LEADER}, {ACTION} that {STRUCTURE}", Action: "Info", Country: "US", Language: "en", IsReadOnly: true, TargetId: "234567", IsFallback: false},

		// {Token: "IA.AR.ARINVOICE.VIEW.LABEL.RECORDTYPE", Text: "Գրանցման տեսակը", Action: "View", Country: "AM", Language: "am", ReadOnly: true},
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

		// Iterate over all the hardcoded records and insert into the database
		for index := range textInfo {
			if "" == textInfo[index].SourceId {
				// If there is no source Id use the value from the config file
				textInfo[index].SourceId = config.ServiceOwnerSourceId
			}
			db.Create(&textInfo[index])
		}

		sql := `CREATE or replace function get_count() returns numeric
	language plpgsql
	as
	$$
	declare
		r_count numeric;
	begin
		SELECT count(id) into r_count FROM text_info;
		return r_count;
	end;
	$$;
	`

		err := db.Exec(sql)
		if err.Error != nil {
			panic(err)
		}

	}

	return db
}
