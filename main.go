package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	cache "github.com/sm-playground/go-text-poc/cache_client"
	cnf "github.com/sm-playground/go-text-poc/config"
	"log"

	"github.com/sm-playground/go-text-poc/db"
	r "github.com/sm-playground/go-text-poc/router"
)

var dbClient *gorm.DB

// main the application entry point
func main() {

	useCache := cnf.GetInstance().Get().Cache.UseCache
	if useCache {
		cacheClient, er := cache.GetCacheClient()
		if er != nil {
			panic("failed to connect to in-memory store")

		} else {
			if cacheClient.Set("hello", "hello world!!!") == nil {
				_ = cacheClient.Delete("hello")
			} else {
				panic("failed to read from in-memory store")
			}
		}
	}

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

	// Initialize router
	r.InitRouter()

}
