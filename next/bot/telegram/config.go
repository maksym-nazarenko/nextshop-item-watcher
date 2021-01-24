package telegram

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

const (
	// TELEGRAM_TOKEN_ENV_VARIABLE holds the name of environment variable where Telegram token resides
	TELEGRAM_TOKEN_ENV_VARIABLE = "NWI_TELEGRAM_TOKEN"
)

// Config holds telegram bot configuration
type Config struct {
	AllowedUsers []subscription.User
}
