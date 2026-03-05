package config

import (
	"os"
	"strconv"
)

// Env returns the value of an environment variable, or the fallback if empty/unset.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// EnvInt returns the value of an environment variable as an int,
// or the fallback if empty/unset or not a valid integer.
func EnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

// EnvBool returns the value of an environment variable as a bool.
// Only "true" and "1" are considered truthy. Everything else returns the fallback.
func EnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v == "true" || v == "1"
}
