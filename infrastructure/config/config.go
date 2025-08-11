package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	ElasticURL  string
	ElasticUser string
	ElasticPass string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPass:      getEnv("DB_PASS", "p@ssw0rd"),
		DBName:      getEnv("DB_NAME", "lms_db"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		ElasticURL:  getEnv("ELASTIC_URL", "http://localhost:9200"),
		ElasticUser: getEnv("ELASTIC_USER", ""),
		ElasticPass: getEnv("ELASTIC_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
