package containers

import "os"

type Config struct {
	Port        string
	Storage     string // "memory" (default) or "sql"
	DatabaseURL string
}

func LoadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		Storage:     getEnv("STORAGE", "memory"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
