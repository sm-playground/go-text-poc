package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

const APP_NAME = "go-text-poc"

type Configuration struct {
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
	Dialect     string
	Port        int
	DBUser      string
	DBPassword  string
	DBName      string
	SSLMode     string
	ShowQueries bool
	DBCP        DBConnectionPool
}

type DBConnectionPool struct {
	MaxIdle   int
	MaxActive int
}

type CacheConfiguration struct {
	UseCache bool
	Dialect  CacheClientType
	Network  string
	IP       string
	Port     int
	Pool     DBConnectionPool
}

type CacheClientType string

const (
	Redis     CacheClientType = "redis"
	Memcached CacheClientType = "memcached"
)

type Config interface {
	Get() Configuration
}

type config struct {
	CF Configuration
}

func (conf *config) Get() Configuration {
	return conf.CF
}

var instance *config

func GetInstance() Config {
	if instance == nil {
		instance = new(config)

		pwd, _ := os.Getwd()

		var tokens []string
		if runtime.GOOS == "windows" {
			tokens = strings.Split(pwd, "\\")
		} else {
			tokens = strings.Split(pwd, "/")
		}

		module := tokens[len(tokens)-1]

		var configFileName string
		if module == APP_NAME {
			configFileName = "config"
		} else {
			configFileName = "config-" + module
		}

		viper.SetConfigName(configFileName)
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		viper.SetConfigType("yml")

		// var configuration c.Configuration
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file, %s", err)
		}

		var configuration Configuration
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

		instance.CF = configuration

	}
	return instance
}
