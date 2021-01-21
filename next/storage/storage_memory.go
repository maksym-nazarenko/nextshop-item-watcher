package storage

import (
	"sync"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

type MemoryStorage struct {
	items     map[string][]subscription.Item
	itemsLock sync.Locker
}

func (m *MemoryStorage) ReadSubscriptions() []subscription.Item {
	ret := make([]subscription.Item, 0, len(m.items))
	for _, item := range m.items {
		ret = append(ret, item...)
	}

	return ret
}

func (m *MemoryStorage) CreateSubscription(item subscription.Item) (bool, error) {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()

	if _, ok := m.items[item.User.ID]; !ok {
		m.items[item.User.ID] = append(m.items[item.User.ID], item)
		return true, nil
	}

	userSubscriptions := m.items[item.User.ID]
	for _, el := range userSubscriptions {
		if item.User.ID == el.User.ID && item.ItemOption.Article == el.ItemOption.Article && item.ItemOption.Number == el.ItemOption.Number {
			return false, nil
		}
	}

	userSubscriptions = append(userSubscriptions, item)
	return true, nil
}

func (m *MemoryStorage) RemoveSubscription(item subscription.Item) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func NewMemoryStorage() subscription.Storage {
	return &MemoryStorage{
		items:     make(map[string][]subscription.Item, 10),
		itemsLock: &sync.Mutex{},
	}
}
