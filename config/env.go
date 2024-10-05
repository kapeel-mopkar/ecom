package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost            string
	AppPort            string
	DBUser             string
	DBPassword         string
	DBAddress          string
	DBName             string
	JWTSecret          string
	JWTExpirationInSec int64
}

const projectDirName = "ecom"

var Envs = initConfig()

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func initConfig() Config {
	loadEnv()
	return Config{
		AppHost:            getEnv("APP_HOST", "http://localhost"),
		AppPort:            getEnv("APP_PORT", "8180"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", "root"),
		DBAddress:          fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		DBName:             getEnv("DB_NAME", "ecom"),
		JWTSecret:          getEnv("JWT_SECRET", "ecom1.0"),
		JWTExpirationInSec: getEnvAsInt("JWT_EXPIRY_SEC", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
