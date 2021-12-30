package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string, fallback string) string {
	err := godotenv.Load(".env")

	if err != nil {
		return fallback
	}

	val := os.Getenv(key)

	if len(val) == 0 {
		return fallback
	}

	return val
}
