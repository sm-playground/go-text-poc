package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sm-playground/go-text-poc/redisClient"
	"log"

	c "github.com/sm-playground/go-text-poc/config"
	"github.com/sm-playground/go-text-poc/db"
	r "github.com/sm-playground/go-text-poc/router"
)

var dbClient *gorm.DB
var config c.Configurations

func main() {

	// Read configuration parameters
	config = c.LoadConfig()

	// Initialize the database and populate with the sample data
	dbClient = db.InitDatabase(config)

	defer func() {
		fmt.Println("\nClose DB connection")
		if err := dbClient.Close(); err != nil {
			log.Printf("ERROR!!! - failed to close DB connection")
		}
	}()

	// Initialize redis connection pool
	redisClient.InitCache(config)

	err := redisClient.Set("hello", "hello world")

	hello, err := redisClient.Get("hello")

	fmt.Println(hello, err)

	r.InitRouter(config, dbClient)

}
