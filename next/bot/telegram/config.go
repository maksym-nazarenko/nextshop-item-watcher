package telegram

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

// Config holds telegram bot configuration
type Config struct {
	AllowedUsers []subscription.User
	Token        string
}
