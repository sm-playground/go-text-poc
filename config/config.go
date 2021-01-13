package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server       ServerConfiguration
	Database     DatabaseConfiguration
	Cache        CacheConfiguration
	InitDatabase bool
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
	DBCP    DBConnectionPool
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

	return configuration
}
