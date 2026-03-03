package config

import (
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type Config struct {
	Server ServerConfig
	Db     DbConfig
	Ai     AiConfig
	// 系統常數 (例如限制、超時)
	MaxChatRetries int //
	Bot            BotConfig
}

type ServerConfig struct {
	Env  string // local, prod
	Port string
}

type AiConfig struct {
	ActiveAI     string
	GeminiAPIKey string
	GeminiModel  string
	OpenAIAPIKey string
}

type BotConfig struct {
	TelegramToken string
}

type DbConfig struct {
	Username        string
	Password        string `json:"-"`
	DbHost          string
	DbPort          int
	DbName          string
	ConnMaxLifetime time.Duration
	LogMode         logger.LogLevel
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "local")
	viper.SetDefault("ACTIVE_AI", "gemini")
	viper.SetDefault("MAX_CHAT_RETRIES", 3)

	// 如果本地有 .env 則讀取
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	viper.AutomaticEnv() // 優先讀取系統環境變數 (Heroku 模式)

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Env:  viper.GetString("ENV"),
		},
		Ai: AiConfig{
			ActiveAI:     viper.GetString("ACTIVE_AI"),
			GeminiAPIKey: viper.GetString("GEMINI_API_KEY"),
			GeminiModel:  viper.GetString("GEMINI_MODEL"),
			OpenAIAPIKey: viper.GetString("OPENAI_API_KEY"),
		},
		Bot: BotConfig{
			TelegramToken: viper.GetString("TELEGRAM_TOKEN"),
		},
		MaxChatRetries: viper.GetInt("MAX_CHAT_RETRIES"),
	}, nil
}
