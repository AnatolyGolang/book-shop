package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	DSN            string
	LogLevel       string
	MigrationsPath string
	HttpPort       string
	HttpHost       string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load("config/local.env")
	if err != nil {
		return Config{}, fmt.Errorf(".env file not found")
	}

	env, err := downloadString("ENV")
	if err != nil {
		return Config{}, fmt.Errorf("can not download env")
	}

	dsn, err := downloadString("DSN")
	if err != nil {
		return Config{}, fmt.Errorf("can not download DB_URL")
	}

	httpPort, err := downloadString("HTTP_PORT")
	if err != nil {
		return Config{}, fmt.Errorf("can not download HTTP_PORT")
	}

	httpHost, err := downloadString("HTTP_HOST")
	if err != nil {
		return Config{}, fmt.Errorf("can not download HTTP_HOST")
	}

	logLevel, err := downloadString("LOG_LEVEL")
	if err != nil {
		return Config{}, fmt.Errorf("can not download LOG_LEVEL")
	}

	migrationsPath, err := downloadString("MIGRATIONS_PATH")
	if err != nil {
		return Config{}, fmt.Errorf("can not download MIGRATIONS_PATH")
	}

	return Config{
		Environment:    env,
		DSN:            dsn,
		HttpPort:       httpPort,
		HttpHost:       httpHost,
		LogLevel:       logLevel,
		MigrationsPath: migrationsPath,
	}, nil
}

func downloadString(key string) (string, error) {
	if val, exists := os.LookupEnv(key); exists {
		return val, nil
	}
	return "", fmt.Errorf("not found value by key %v", key)
}
