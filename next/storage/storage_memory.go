package storage

import (
	"sync"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

// MemoryStorage describes in-memory storage for subscriptions
type MemoryStorage struct {
	items     map[string][]subscription.Item
	itemsLock sync.Locker
}

// ReadSubscriptions reads all subscriptions from subscription storage
func (m *MemoryStorage) ReadSubscriptions() []subscription.Item {
	ret := make([]subscription.Item, 0, len(m.items))
	for _, item := range m.items {
		ret = append(ret, item...)
	}

	return ret
}

// CreateSubscription creates new subscription from subscription storage
func (m *MemoryStorage) CreateSubscription(item *subscription.Item) (bool, error) {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()

	if _, ok := m.items[item.User.ID]; !ok {
		m.items[item.User.ID] = append(m.items[item.User.ID], *item)
		return true, nil
	}

	userSubscriptions := m.items[item.User.ID]
	for _, el := range userSubscriptions {
		if item.User.ID == el.User.ID && item.ShopItem.Article == el.ShopItem.Article && item.ShopItem.SizeID == el.ShopItem.SizeID {
			return false, nil
		}
	}

	m.items[item.User.ID] = append(userSubscriptions, *item)

	return true, nil
}

// RemoveSubscription removes subscription from subscription storage
func (m *MemoryStorage) RemoveSubscription(item *subscription.Item) (bool, error) {
	panic("not implemented") // TODO: Implement
}

// NewMemoryStorage constructs new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items:     make(map[string][]subscription.Item, 10),
		itemsLock: &sync.Mutex{},
	}
}
