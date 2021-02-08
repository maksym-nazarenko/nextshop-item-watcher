package mediator

import (
	"fmt"
	"log"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

// SubscriptionStorage describes storage-related actions
type SubscriptionStorage interface {
	CreateSubscription(subscription.Item) (bool, error)
	DisableSubscription(subscription.Item) error
	EnableSubscription(subscription.Item) error
	ReadAllSubscriptions() ([]subscription.Item, error)
	ReadSubscriptions() ([]subscription.Item, error)
	ReadSubscriptionsByShopItem(shop.Item) ([]subscription.Item, error)
	ReadUserAllSubscriptions(subscription.User) ([]subscription.Item, error)
	ReadUserSubscriptions(subscription.User) ([]subscription.Item, error)
	RemoveSubscription(subscription.Item) (bool, error)
}

// SubscriptionMediator de-couples different components of the system
type SubscriptionMediator struct {
	StorageBackend SubscriptionStorage

	watcher       watch.Watcher
	inStockItemCh chan subscription.Item
	httpClient    *next.Client
}

// ReadSubscriptions reads all subscriptions
func (m *SubscriptionMediator) ReadSubscriptions() ([]subscription.Item, error) {
	return m.StorageBackend.ReadSubscriptions()
}

// CreateSubscription creates new subscription in system
func (m *SubscriptionMediator) CreateSubscription(item subscription.Item) (bool, error) {
	extendedOptions, err := m.httpClient.GetItemExtendedOption(item.ShopItem.Article)

	if err != nil {
		log.Println("[ERROR] Could not enrich subscription item with extra data: " + err.Error())
	}

	item.ShopItem.Description = extendedOptions.Description
	option, ok := m.httpClient.FindOptionBySize(extendedOptions.Options, item.ShopItem.SizeID)
	if ok {
		item.ShopItem.SizeString = option.Name
	}

	url, err := m.httpClient.GetItemURLByArticle(item.ShopItem.Article)
	if err != nil {
		log.Println("[ERROR] Could not fetch item URL: " + err.Error())
	}

	item.ShopItem.URL = url

	ok, err = m.StorageBackend.CreateSubscription(item)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	if err = m.watcher.AddItem(&item.ShopItem); err != nil {
		return false, err
	}

	return true, nil
}

// RemoveSubscription removes subscription from system
func (m *SubscriptionMediator) RemoveSubscription(item subscription.Item) (bool, error) {
	return m.StorageBackend.RemoveSubscription(item)
}

func (m *SubscriptionMediator) FetchSizeIDs(article string) ([]shop.ItemOption, error) {
	items, err := m.httpClient.GetOptionsByArticle(article)

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Start begins the main loop
func (m *SubscriptionMediator) Start() {
	var item subscription.Item
	var err error

	for inStockItem := range m.watcher.InStockChan() {
		log.Printf("[DEBUG] item appeared in stock: %v", inStockItem)
		if item, err = m.findItemByShopItem(inStockItem); err != nil {
			log.Printf("[ERROR] %s\n", err.Error())
			continue
		}

		m.inStockItemCh <- item

		m.watcher.RemoveItem(shop.Item{Article: item.ShopItem.Article, SizeID: item.ShopItem.SizeID})
		log.Printf("[DEBUG] mediator: disabling subscription %v\n", item)
		if err := m.StorageBackend.DisableSubscription(item); err != nil {
			log.Printf("[ERROR] mediator: could not disable subscription: %s\n", err.Error())
		}
	}
}

func (m *SubscriptionMediator) Stop() {
	log.Println("[INFO] Stopping mediator")
}

func (m *SubscriptionMediator) InStockItemCh() <-chan subscription.Item {
	return m.inStockItemCh
}

func (m *SubscriptionMediator) findItemByShopItem(item shop.Item) (subscription.Item, error) {
	subscriptions, err := m.ReadSubscriptions()
	if err != nil {
		return subscription.Item{}, err
	}

	for _, it := range subscriptions {
		if it.ShopItem.Article == item.Article && it.ShopItem.SizeID == item.SizeID {
			return it, nil
		}
	}

	return subscription.Item{}, fmt.Errorf("no such subscription item found: %v", item)
}

// New instantiates SubscriptionMediator object
func New(storageBackend SubscriptionStorage, watcher watch.Watcher, httpClient *next.Client) *SubscriptionMediator {
	return &SubscriptionMediator{
		StorageBackend: storageBackend,
		inStockItemCh:  make(chan subscription.Item, 10),
		watcher:        watcher,
		httpClient:     httpClient,
	}
}
