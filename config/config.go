package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string
	GRPCServer GRPCServer
	Database   Database
	Prometheus Prometheus
}

type GRPCServer struct {
	Address string
	Port    int
}

type Database struct {
	Username       string
	Password       string
	Host           string
	Port           string
	DbName         string
	MigrationsPath string
}

type Prometheus struct {
	Address string
	Port    int
}

func MustLoad() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("env", "dev")

	viper.SetDefault("grpc_server.address", "0.0.0.0")
	viper.SetDefault("grpc_server.port", 50051)

	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "admin")
	viper.SetDefault("database.host", "user-db")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.db_name", "userservice")
	viper.SetDefault("database.migrations_path", "migrations")

	viper.SetDefault("prometheus.address", "0.0.0.0")
	viper.SetDefault("prometheus.port", 9101)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		os.Exit(1)
	}

	config := &Config{
		Env: viper.GetString("env"),
		GRPCServer: GRPCServer{
			Address: viper.GetString("grpc_server.address"),
			Port:    viper.GetInt("grpc_server.port"),
		},
		Database: Database{
			Username:       viper.GetString("database.username"),
			Password:       viper.GetString("database.password"),
			Host:           viper.GetString("database.host"),
			Port:           viper.GetString("database.port"),
			DbName:         viper.GetString("database.db_name"),
			MigrationsPath: viper.GetString("database.migrations_path"),
		},
		Prometheus: Prometheus{
			Address: viper.GetString("prometheus.address"),
			Port:    viper.GetInt("prometheus.port"),
		},
	}

	return config
}
