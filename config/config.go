package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string
	HTTPServer HTTPServer
	Database   Database
	JWT        JWT
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Database struct {
	Username       string
	Password       string
	Host           string
	Port           string
	DbName         string
	MigrationsPath string
}

type JWT struct {
	Secret           string
	AccessExpiresAt  time.Duration
	RefreshExpiresAt time.Duration
}

func MustLoad() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("env", "dev")

	viper.SetDefault("http_server.address", "0.0.0.0:8000")
	viper.SetDefault("http_server.timeout", "4s")
	viper.SetDefault("http_server.idle_timeout", "60s")

	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "admin")
	viper.SetDefault("database.host", "db")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.db_name", "userservice")
	viper.SetDefault("database.migrations_path", "migrations")

	viper.SetDefault("jwt.secret", "my-secret")
	viper.SetDefault("jwt.access_expires_at", "1m")
	viper.SetDefault("jwt.refresh_expires_at", "5m")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		os.Exit(1)
	}

	config := &Config{
		Env: viper.GetString("env"),
		HTTPServer: HTTPServer{
			Address:     viper.GetString("http_server.address"),
			Timeout:     viper.GetDuration("http_server.timeout"),
			IdleTimeout: viper.GetDuration("http_server.idle_timeout"),
		},
		Database: Database{
			Username:       viper.GetString("database.username"),
			Password:       viper.GetString("database.password"),
			Host:           viper.GetString("database.host"),
			Port:           viper.GetString("database.port"),
			DbName:         viper.GetString("database.db_name"),
			MigrationsPath: viper.GetString("database.migrations_path"),
		},
		JWT: JWT{
			Secret:           viper.GetString("jwt.secret"),
			AccessExpiresAt:  viper.GetDuration("jwt.access_expires_at"),
			RefreshExpiresAt: viper.GetDuration("jwt.refresh_expires_at"),
		},
	}

	return config
}
