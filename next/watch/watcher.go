package watch

import (
	"log"
	"sync"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/robfig/cron/v3"
)

// Watcher interface to be implemented by different watchers
type Watcher interface {
	AddItem(*shop.Item) error
	InStockChan() <-chan shop.Item
	RemoveItem(shop.Item)
	Start() error
	Stop() error
}

// ItemWatcher holds information about items to watch after
type ItemWatcher struct {
	Client         *next.Client
	UpdateInterval time.Duration
	cron           *cron.Cron
	items          []*shop.Item
	itemsLock      sync.Locker
	inStockChan    chan shop.Item
}

// Start begins watcher's loop of checks
func (w *ItemWatcher) Start() error {
	w.cron.Start()

	return nil
}

// Stop terminates periodic checking
func (w *ItemWatcher) Stop() error {
	defer close(w.inStockChan)
	w.cron.Stop()

	return nil
}

// Run satisfies cron.Job interface
func (w *ItemWatcher) Run() {
	w.onTimer()
}

func (w *ItemWatcher) onTimer() {
	log.Println("ItemWatcher timer fired")
	for _, item := range w.items {
		go func(item shop.Item) {
			shopItemInfo, err := w.Client.GetItemInfo(item)
			if err != nil {
				log.Println("[ERROR] + " + err.Error())
				return
			}

			w.processInStockItems(shopItemInfo)
		}(*item)
	}
}

// AddItem add given item to the list of watched items
func (w *ItemWatcher) AddItem(item *shop.Item) error {
	w.items = append(w.items, item)

	return nil
}

// RemoveItem removes watching item from the list so it will not be processed next time when cron fires
func (w *ItemWatcher) RemoveItem(item shop.Item) {
	w.itemsLock.Lock()
	defer w.itemsLock.Unlock()
	for index, it := range w.items {
		if item.Article == it.Article && item.SizeID == item.SizeID {
			w.items = append(w.items[:index], w.items[index+1:]...)
			return
		}
	}
}

// InStockChan returns channel where InStock items will appear
func (w ItemWatcher) InStockChan() <-chan shop.Item {
	return w.inStockChan
}

func (w *ItemWatcher) processInStockItems(items ...shop.ItemOption) {
	w.itemsLock.Lock()
	defer w.itemsLock.Unlock()
	for _, item := range items {
		if item.StockStatusString != shop.ItemStatusInStock {
			continue
		}
		w.inStockChan <- shop.Item{Article: item.Article, SizeID: item.Number}
	}
}

// New constructs new ItemWatcher instance
func New(client *next.Client, config *Config) *ItemWatcher {
	// TODO: add TZ support
	watcher := ItemWatcher{
		Client:         client,
		UpdateInterval: config.UpdateInterval,
		cron:           cron.New(),
		itemsLock:      &sync.Mutex{},
	}

	watcher.inStockChan = make(chan shop.Item, 20)
	interval := "@every " + watcher.UpdateInterval.String()
	watcher.cron.AddJob(interval, &watcher)

	return &watcher
}
