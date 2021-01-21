package mediator

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

type Mediator interface {
	subscription.Reader
	subscription.Writer
}

type SubscriptionMediator struct {
	StorageBackend subscription.Storage

	watchers                    map[string]watch.Watcher
	newSubscriptionItemChan     chan subscription.Item
	subscriptionItemCreatedChan chan subscription.Item
}

func (m *SubscriptionMediator) ReadSubscriptions() []subscription.Item {
	return m.StorageBackend.ReadSubscriptions()
}

func (m *SubscriptionMediator) CreateSubscription(item subscription.Item) (bool, error) {
	return m.CreateSubscription(item)
}

func (m *SubscriptionMediator) RemoveSubscription(item subscription.Item) (bool, error) {
	return m.RemoveSubscription(item)
}

func New(storageBackend subscription.Storage) Mediator {
	return &SubscriptionMediator{
		StorageBackend:              storageBackend,
		watchers:                    make(map[string]watch.Watcher),
		newSubscriptionItemChan:     make(chan subscription.Item),
		subscriptionItemCreatedChan: make(chan subscription.Item),
	}
}
