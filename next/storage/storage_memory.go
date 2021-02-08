package storage

import (
	"errors"
	"sync"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

// MemoryStorage describes in-memory storage for subscriptions
type MemoryStorage struct {
	items     map[string][]*subscription.Item
	itemsLock sync.Locker
}

// ReadSubscriptions reads all subscriptions from subscription storage
func (m *MemoryStorage) ReadSubscriptions() ([]subscription.Item, error) {
	ret := make([]subscription.Item, 0, len(m.items))
	for _, items := range m.items {
		for _, item := range items {
			if item.Active {
				ret = append(ret, *item)
			}
		}
	}

	return ret, nil
}

// CreateSubscription creates new subscription from subscription storage
func (m *MemoryStorage) CreateSubscription(item subscription.Item) (bool, error) {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()

	if _, ok := m.items[item.User.ID]; !ok {
		m.items[item.User.ID] = append(m.items[item.User.ID], &item)
		return true, nil
	}

	userSubscriptions := m.items[item.User.ID]
	for _, el := range userSubscriptions {
		if item.User.ID == el.User.ID && item.ShopItem.Article == el.ShopItem.Article && item.ShopItem.SizeID == el.ShopItem.SizeID {
			return false, nil
		}
	}

	m.items[item.User.ID] = append(userSubscriptions, &item)

	return true, nil
}

func (m *MemoryStorage) DisableSubscription(item subscription.Item) error {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()
	userItem, err := m.findUserItem(item)

	if err != nil {
		return err
	}

	userItem.Active = false

	return nil
}

func (m *MemoryStorage) EnableSubscription(item subscription.Item) error {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()
	userItem, err := m.findUserItem(item)

	if err != nil {
		return err
	}

	userItem.Active = true

	return nil
}

func (m *MemoryStorage) ReadAllSubscriptions() ([]subscription.Item, error) {
	ret := make([]subscription.Item, 0, len(m.items))
	for _, items := range m.items {
		for _, item := range items {
			ret = append(ret, *item)
		}
	}

	return ret, nil
}

func (m *MemoryStorage) ReadSubscriptionsByShopItem(item shop.Item) ([]subscription.Item, error) {
	ret := make([]subscription.Item, 0, len(m.items))
	for _, items := range m.items {
		for _, userItem := range items {
			if userItem.Active && userItem.ShopItem.Article == item.Article && userItem.ShopItem.SizeID == item.SizeID {
				ret = append(ret, *userItem)
			}
		}
	}

	return ret, nil
}

func (m *MemoryStorage) ReadUserSubscriptions(user subscription.User) ([]subscription.Item, error) {
	userItems, ok := m.items[user.ID]
	if !ok {
		return nil, errors.New("No such user found")
	}

	ret := make([]subscription.Item, 0, len(userItems))
	for _, item := range userItems {
		if item.Active {
			ret = append(ret, *item)
		}
	}

	return ret, nil
}

func (m *MemoryStorage) ReadUserAllSubscriptions(user subscription.User) ([]subscription.Item, error) {
	userItems, ok := m.items[user.ID]
	if !ok {
		return nil, errors.New("No such user found")
	}

	ret := make([]subscription.Item, 0, len(userItems))
	for _, item := range userItems {
		ret = append(ret, *item)
	}

	return ret, nil
}

func (m *MemoryStorage) RemoveSubscription(item subscription.Item) (bool, error) {
	m.itemsLock.Lock()
	defer m.itemsLock.Unlock()

	userSubscriptions, ok := m.items[item.User.ID]
	if !ok {
		return false, nil
	}

	for index, userItem := range userSubscriptions {
		if userItem.ShopItem.Article == item.ShopItem.Article && userItem.ShopItem.SizeID == item.ShopItem.SizeID {
			userSubscriptions[index] = userSubscriptions[len(userSubscriptions)-1]
			m.items[item.User.ID] = userSubscriptions[:len(userSubscriptions)-1]

			return true, nil
		}
	}

	return false, nil
}

func (m *MemoryStorage) findUserItem(item subscription.Item) (*subscription.Item, error) {

	userSubscriptions, ok := m.items[item.User.ID]
	if !ok {
		return nil, errors.New("No user for subscription found")
	}

	for _, userItem := range userSubscriptions {
		if userItem.ShopItem.Article == item.ShopItem.Article && userItem.ShopItem.SizeID == item.ShopItem.SizeID {
			return userItem, nil
		}
	}

	return nil, errors.New("No subscription found")
}

// NewMemoryStorage constructs new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items:     make(map[string][]*subscription.Item),
		itemsLock: &sync.Mutex{},
	}
}
