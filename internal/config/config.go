package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	Yandex      YandexConfig
}

type YandexConfig struct {
	APIKey   string
	Model    string
	BaseURL  string
	FolderID string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getenv("PORT", "8080")
	return &Config{
		Port:        port,
		Env:         getenv("NODE_ENV", "development"),
		DatabaseURL: getenv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/acc_db?sslmode=disable"),
		Yandex: YandexConfig{
			APIKey:   os.Getenv("YANDEX_AI_API_KEY"),
			Model:    os.Getenv("YANDEX_AI_MODEL"),
			BaseURL:  os.Getenv("YANDEX_AI_BASE_URL"),
			FolderID: os.Getenv("YANDEX_FOLDER_ID"),
		},
	}, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ParsePort(port string) int {
	p, err := strconv.Atoi(port)
	if err != nil || p <= 0 {
		return 8080
	}
	return p
}
