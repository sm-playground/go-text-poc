package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Configurations struct {
	Server               ServerConfiguration
	Database             DatabaseConfiguration
	Cache                CacheConfiguration
	InitDatabase         bool
	DefaultLocale        string
	DefaultCountry       string
	DefaultLanguage      string
	ServiceOwnerSourceId string
}

type ServerConfiguration struct {
	IP   string
	Port int
}

type DatabaseConfiguration struct {
	Dialect    string
	Port       int
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	DBCP       DBConnectionPool
}

type DBConnectionPool struct {
	MaxIdle   int
	MaxActive int
}

type CacheConfiguration struct {
	Network string
	IP      string
	Port    int
	Pool    DBConnectionPool
}

// Loads the application configuration parameters from the yaml file
func LoadConfig() (configuration Configurations) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	// var configuration c.Configurations
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
	}

	l := strings.Split(configuration.DefaultLocale, "-")

	if len(l) != 2 {
		// the locale is expected in the format of "en-US"
		// Ignore all other cases
		panic("incorrect format for locale: " + fmt.Sprintf("%+v", configuration.DefaultLocale))
	}
	configuration.DefaultLanguage = strings.TrimSpace(l[0])
	configuration.DefaultCountry = strings.TrimSpace(l[1])

	return configuration
}
