package mediator

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

// SubscriptionStorage describes storage-related actions
type SubscriptionStorage interface {
	ReadSubscriptions() []subscription.Item
	CreateSubscription(subscription.Item) (bool, error)
	RemoveSubscription(subscription.Item) (bool, error)
}

// SubscriptionMediator de-couples different components of the system
type SubscriptionMediator struct {
	StorageBackend SubscriptionStorage

	watchers                    map[string]watch.Watcher
	newSubscriptionItemChan     chan subscription.Item
	subscriptionItemCreatedChan chan subscription.Item
}

// ReadSubscriptions reads all subscriptions
func (m *SubscriptionMediator) ReadSubscriptions() []subscription.Item {
	return m.StorageBackend.ReadSubscriptions()
}

// CreateSubscription creates new subscription in system
func (m *SubscriptionMediator) CreateSubscription(item subscription.Item) (bool, error) {
	return m.CreateSubscription(item)
}

// RemoveSubscription remoces subscription from system
func (m *SubscriptionMediator) RemoveSubscription(item subscription.Item) (bool, error) {
	return m.RemoveSubscription(item)
}

// New instantiates SubscriptionMediator object
func New(storageBackend subscription.Storage) SubscriptionMediator {
	return SubscriptionMediator{
		StorageBackend:              storageBackend,
		watchers:                    make(map[string]watch.Watcher),
		newSubscriptionItemChan:     make(chan subscription.Item),
		subscriptionItemCreatedChan: make(chan subscription.Item),
	}
}
