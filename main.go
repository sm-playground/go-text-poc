package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sm-playground/go-text-poc/redisClient"
	"log"

	"github.com/sm-playground/go-text-poc/db"
	r "github.com/sm-playground/go-text-poc/router"
)

var dbClient *gorm.DB

func main() {

	// Initialize the database and populate with the sample data
	dbClient = db.GetConnection()

	// var textInfo []m.TextInfo
	// dbClient.Raw("select * from text_info").Scan(textInfo)

	defer func() {
		fmt.Println("\nClose DB connection")
		if err := dbClient.Close(); err != nil {
			log.Printf("ERROR!!! - failed to close DB connection")
		}
	}()

	// Initialize redis connection pool
	redisClient.InitCache()

	err := redisClient.Set("hello", "hello world")

	hello, err := redisClient.Get("hello")

	fmt.Println(hello, err)

	r.InitRouter()

}
