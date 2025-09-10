package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	Server   ServerConfig   `mapstructure:"server"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	AI       AIConfig       `mapstructure:"ai"`
	Bot      BotConfig      `mapstructure:"bot"`
}

// TelegramConfig holds Telegram bot configuration
type TelegramConfig struct {
	Token          string `mapstructure:"token"`
	APIBaseURL     string `mapstructure:"api_base_url"`
	Debug          bool   `mapstructure:"debug"`
	UpdatesTimeout int    `mapstructure:"updates_timeout"`
	WebhookEnable  bool   `mapstructure:"webhook_enable"`
	WebhookURL     string `mapstructure:"webhook_url"`
	WebhookDomain  string `mapstructure:"webhook_domain"`
	WebhookPath    string `mapstructure:"webhook_path"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address string `mapstructure:"address"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	URL        string `mapstructure:"url"`
	Model      string `mapstructure:"model"`
	APIKey     string `mapstructure:"api_key"`
	Prompt     string `mapstructure:"prompt"`
	PromptFile string `mapstructure:"prompt_file"`
}

// BotConfig holds bot messages and behavior configuration
type BotConfig struct {
	StartMessage          string `mapstructure:"start_message"`
	HelpMessage           string `mapstructure:"help_message"`
	UnknownCommandMessage string `mapstructure:"unknown_command_message"`
	ErrorMessage          string `mapstructure:"error_message"`
	EmptyMessage          string `mapstructure:"empty_message"`
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, ignore error
		_ = err
	}

	// Set default values
	viper.SetDefault("telegram.debug", false)
	viper.SetDefault("telegram.updates_timeout", 30)
	viper.SetDefault("telegram.webhook_enable", false)
	viper.SetDefault("telegram.webhook_path", "/webhook")
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("ai.model", "gpt-3.5-turbo")

	// Bot message defaults
	viper.SetDefault("bot.start_message", "🤖 Привет! Я универсальный AI-ассистент.\n\n💡 Просто отправьте мне сообщение, и я помогу вам с любыми вопросами!\n\nИспользуйте /help для получения дополнительной информации.")
	viper.SetDefault("bot.help_message", "📚 Помощь по использованию AI-ассистента:\n\n💬 **Любое сообщение** → Получите умный ответ:\n• Ответы на вопросы\n• Помощь с задачами\n• Объяснения и советы\n• Творческие идеи\n\n🔧 **Доступные команды:**\n• /start - Начать работу с ботом\n• /help - Показать эту справку\n\n💡 Просто отправьте текст - я сразу помогу!")
	viper.SetDefault("bot.unknown_command_message", "❓ Неизвестная команда. Используйте /help для получения информации о возможностях бота.")
	viper.SetDefault("bot.error_message", "Извините, произошла ошибка при обработке вашего сообщения. Попробуйте еще раз.")
	viper.SetDefault("bot.empty_message", "Пожалуйста, отправьте текстовое сообщение.")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind environment variables to viper keys
	_ = viper.BindEnv("telegram.token", "TELEGRAM_TOKEN")
	_ = viper.BindEnv("telegram.api_base_url", "TELEGRAM_API_BASE_URL")
	_ = viper.BindEnv("telegram.debug", "TELEGRAM_DEBUG")
	_ = viper.BindEnv("telegram.updates_timeout", "TELEGRAM_UPDATES_TIMEOUT")
	_ = viper.BindEnv("telegram.webhook_enable", "TELEGRAM_WEBHOOK_ENABLE")
	_ = viper.BindEnv("telegram.webhook_url", "TELEGRAM_WEBHOOK_URL")
	_ = viper.BindEnv("telegram.webhook_domain", "TELEGRAM_WEBHOOK_DOMAIN")
	_ = viper.BindEnv("telegram.webhook_path", "TELEGRAM_WEBHOOK_PATH")
	_ = viper.BindEnv("server.address", "SERVER_ADDRESS")
	_ = viper.BindEnv("logging.level", "LOG_LEVEL")
	_ = viper.BindEnv("ai.url", "AI_URL")
	_ = viper.BindEnv("ai.model", "AI_MODEL")
	_ = viper.BindEnv("ai.api_key", "AI_API_KEY")
	_ = viper.BindEnv("ai.prompt", "AI_PROMPT")
	_ = viper.BindEnv("ai.prompt_file", "AI_PROMPT_FILE")
	_ = viper.BindEnv("bot.start_message", "BOT_START_MESSAGE")
	_ = viper.BindEnv("bot.help_message", "BOT_HELP_MESSAGE")
	_ = viper.BindEnv("bot.unknown_command_message", "BOT_UNKNOWN_COMMAND_MESSAGE")
	_ = viper.BindEnv("bot.error_message", "BOT_ERROR_MESSAGE")
	_ = viper.BindEnv("bot.empty_message", "BOT_EMPTY_MESSAGE")

	// Set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use env vars and defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Load prompt from file if specified
	if config.AI.PromptFile != "" {
		promptFromFile, err := loadPromptFromFile(config.AI.PromptFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load prompt from file: %w", err)
		}
		config.AI.Prompt = promptFromFile
	}

	// Validate required fields
	if config.Telegram.Token == "" {
		return nil, fmt.Errorf("telegram token is required")
	}
	if config.AI.URL == "" {
		return nil, fmt.Errorf("ai url is required")
	}
	if config.AI.APIKey == "" {
		return nil, fmt.Errorf("ai api key is required")
	}
	if config.AI.Prompt == "" {
		return nil, fmt.Errorf("ai prompt is required (either AI_PROMPT or AI_PROMPT_FILE must be set)")
	}

	return &config, nil
}

// loadPromptFromFile loads prompt content from a file
func loadPromptFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", filePath, err)
	}
	return strings.TrimSpace(string(content)), nil
}

// GetEnv returns environment variable value or default
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
