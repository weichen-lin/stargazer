package infrastructure

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	ClerkSecret     string
	BackgroundToken string
	DatabaseURL     string
	ResendApiKey    string
	RedisURL        string
}

func NewConfig() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")

	// 只有在非 production 環境下才載入 .env 文件
	if appEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	clerkSecret := os.Getenv("CLERK_SECRET_KEY")
	backgroundToken := os.Getenv("BACKGROUND_SERVICE_TOKEN")
	databaseURL := os.Getenv("DATABASE_URL")
	resendApiKey := os.Getenv("RESEND_API_KEY")
	redisURL := os.Getenv("REDIS_URL")

	if clerkSecret == "" ||
		backgroundToken == "" ||
		databaseURL == "" ||
		resendApiKey == "" || redisURL == "" {
		return nil,
			fmt.Errorf("missing required environment variables %s, %s, %s, %s, %s, %s, %s", appEnv, clerkSecret, backgroundToken, databaseURL, resendApiKey, redisURL)

	}

	config := &Config{
		AppEnv:          appEnv,
		ClerkSecret:     clerkSecret,
		BackgroundToken: backgroundToken,
		DatabaseURL:     databaseURL,
		ResendApiKey:    resendApiKey,
		RedisURL:        redisURL,
	}

	return config, nil
}
