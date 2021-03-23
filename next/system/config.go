package system

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/bot/telegram"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/storage"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

// Config holds configuration for the whole system
// basically, it wraps all other configurations
type Config struct {
	HTTP    HTTPConfig
	Watch   watch.Config
	Bot     telegram.Config
	Storage storage.Config
}

type HTTPConfig struct {
	Client next.Config
}

type StorageConfig struct {
	Type    string
	Options map[string]interface{}
}
