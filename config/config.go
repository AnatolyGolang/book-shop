package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	DSN            string
	LogLevel       string
	MigrationsPath string
	HttpPort       int
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

	httpPort, err := downloadInt("HTTP_PORT")
	if err != nil {
		return Config{}, fmt.Errorf("can not download HTTP_PORT")
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

func downloadInt(key string) (int, error) {
	if val, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal, nil
		}
	}
	return 0, fmt.Errorf("not found value by key %v", key)
}
