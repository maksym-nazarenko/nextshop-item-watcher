package watcher

import (
	"log"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/robfig/cron/v3"
)

// Watcher interface to be implemented by different watchers
type Watcher interface {
	AddItem(next.ShopItem) error
	AddHandlers(...ItemOptionHandler)
	AddHandlerFuncs(...handleFuncType)
	Start() error
	Stop() error
}

// ItemWatcher holds information about items to watch after
type ItemWatcher struct {
	Client         *next.Client
	UpdateInterval time.Duration
	cron           *cron.Cron
	items          []next.ShopItem
	handlers       []ItemOptionHandler
}

// AddHandlers adds one or more handlers to be notified when the item becomes available
func (w *ItemWatcher) AddHandlers(handlers ...ItemOptionHandler) {
	w.handlers = append(w.handlers, handlers...)
}

// AddHandlerFuncs adds one or more handler functions to be called when the item becomes available
func (w *ItemWatcher) AddHandlerFuncs(handleFuncs ...handleFuncType) {
	w.handlers = make([]ItemOptionHandler, 0, len(handleFuncs))
	for _, f := range handleFuncs {
		w.handlers = append(w.handlers, HandleFuncWrapper{handleFunc: f})
	}
}

// Start begins watcher's loop of checks
func (w *ItemWatcher) Start() error {
	// TODO: this should be added only once
	interval := "@every " + w.UpdateInterval.String()
	w.cron.AddFunc(interval, w.onTimer)

	w.cron.Start()

	return nil
}

// Stop terminates periodic checking
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

			w.processInStockItems(shopItemInfo)
		}(item)
	}
}

// AddItem add given item to the list of watched items
func (w *ItemWatcher) AddItem(item next.ShopItem) error {
	w.items = append(w.items, item)

	return nil
}

func (w *ItemWatcher) processInStockItems(items ...next.ItemOption) {

	inStockItems := make([]next.ItemOption, 0, 10)
	for _, item := range items {
		if item.StockStatusString != next.ItemStatusInStock {
			continue
		}
		inStockItems = append(inStockItems, item)
	}

	if len(inStockItems) < 1 {
		return
	}

	for _, h := range w.handlers {
		h.Handle(inStockItems...)
	}
}

// New constructs new Watcher instance
func New(client *next.Client, config *Config) Watcher {
	// TODO: add TZ support
	watcher := ItemWatcher{Client: client, UpdateInterval: config.UpdateInterval, cron: cron.New()}

	return &watcher
}
