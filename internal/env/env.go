package env

import (
	"os"
	"strconv"
)

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return i
}

func GetString(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
