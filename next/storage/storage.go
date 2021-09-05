package storage

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

type Storage interface {
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
