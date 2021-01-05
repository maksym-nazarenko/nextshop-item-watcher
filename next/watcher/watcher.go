package watcher

import (
	"log"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/robfig/cron/v3"
)

type Watcher interface {
	AddItem(item next.ShopItem) error
	Start() error
	Stop() error
}

type ItemWatcher struct {
	Client         *next.Client
	UpdateInterval time.Duration
	cron           *cron.Cron
	items          []next.ShopItem
}

func (w *ItemWatcher) Start() error {
	// TODO: this should be added only once
	interval := "@every " + w.UpdateInterval.String()
	w.cron.AddFunc(interval, w.onTimer)

	w.cron.Start()

	return nil
}

func (w *ItemWatcher) Stop() error {
	w.cron.Stop()

	return nil
}

func (w *ItemWatcher) onTimer() {
	log.Println("ItemWatcher timer fired")
	for _, item := range w.items {
		go func(item next.ShopItem) {
			shopItemInfo, err := w.Client.GetItemInfo(item)
			if err != nil {
				log.Println("[ERROR] + " + err.Error())
				return
			}

			log.Printf("[%s][%s] %s - %s", item.Article, shopItemInfo.StockStatusString, shopItemInfo.Name, shopItemInfo.Price)

		}(item)
	}
}

func (w *ItemWatcher) AddItem(item next.ShopItem) error {
	w.items = append(w.items, item)

	return nil
}

func New(client *next.Client, config *Config) Watcher {
	// TODO: add TZ support
	watcher := ItemWatcher{Client: client, UpdateInterval: config.UpdateInterval, cron: cron.New()}

	return &watcher
}
