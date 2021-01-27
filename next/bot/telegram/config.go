package telegram

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

const (
	// TelegramTokenEnvVariable holds the name of environment variable where Telegram token resides
	TelegramTokenEnvVariable = "NWI_TELEGRAM_TOKEN"
)

// Config holds telegram bot configuration
type Config struct {
	AllowedUsers []subscription.User
	Token        string
}
