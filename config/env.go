package config

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type PostgresConfig struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	DBName   string
// 	SSLMode  string
// }

// var Envs = initPostgresConfig()

// func initPostgresConfig() PostgresConfig {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	return PostgresConfig{
// 		Host:     getEnv("DB_HOST", "localhost"),
// 		Port:     getEnv("DB_PORT", "5432"),
// 		User:     getEnv("DB_USER", "postgres"),
// 		Password: getEnv("DB_PASSWORD", ""),
// 		DBName:   getEnv("DB_NAME", "mydatabase"),
// 		SSLMode:  getEnv("DB_SSLMODE", "disable"),
// 	}
// }

// func getEnv(key, fallback string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}

// 	return fallback
// }
