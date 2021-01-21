package telegram

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

const (
	TELEGRAM_TOKEN_ENV_VARIABLE = "NWI_TELEGRAM_TOKEN"
)

// Config holds telegram bot configuration
type Config struct {
	// httpClient   *next.Client
	// mediator     *mediator.Mediator
	AllowedUsers []subscription.User
}
