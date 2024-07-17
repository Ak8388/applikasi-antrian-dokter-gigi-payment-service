package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiConfig
	DBConfig
	JwtConfig
}

type ApiConfig struct {
	ApiPort string
}

type DBConfig struct {
	DBUser   string
	DBPass   string
	DBName   string
	DBPort   string
	DBDriver string
}

type JwtConfig struct {
	SecretKey string
}

func (cfg *Config) readConfig() error {

	if err := godotenv.Load(); err != nil {
		return errors.New("failed read config from environtment")
	}

	cfg.DBConfig = DBConfig{
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		DBPort:   os.Getenv("DB_PORT"),
		DBDriver: os.Getenv("DB_DRIVER"),
	}

	cfg.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("API_PORT"),
	}

	cfg.JwtConfig = JwtConfig{
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	if cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" || cfg.DBPort == "" || cfg.DBDriver == "" || cfg.ApiPort == "" {
		return errors.New("all environtments required")
	}

	return nil
}

func Cfg() *Config {
	cfg := &Config{}

	if err := cfg.readConfig(); err != nil {
		panic(err)
	}

	return cfg
}
