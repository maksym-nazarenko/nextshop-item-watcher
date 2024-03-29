package storage

import (
	"testing"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"github.com/stretchr/testify/assert"
)

func TestStorageMemory_addItemToEmptyStorage(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	added, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  1,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)

	assert.Equal(1, len(subscriptions))
}

func TestStorageMemory_addItemAddsSameItemOnlyOnce(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	added, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	added, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)

	assert.False(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))
}

func TestStorageMemory_addItemAddsSecondItemIfDifferent(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	added, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-333",
				SizeID:  11,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	added, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-333",
				SizeID:  12,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_readSubscriptionsFetchesOnlyActive(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	added, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	added, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-5"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	added, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-2"},
			ShopItem: shop.Item{
				Article: "333-444",
				SizeID:  25,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_enableSubscription(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-5"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)

	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	err = strg.EnableSubscription(
		subscription.Item{
			User:     subscription.User{ID: "user-5"},
			ShopItem: shop.Item{Article: "222-123", SizeID: 10},
		},
	)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_readAllSubscriptions(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)

	_, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-5"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)

	assert.NoError(err)

	subscriptions, err := strg.ReadAllSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_readUserSubscription(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)

	_, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  11,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "333-222",
				SizeID:  15,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-3"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	subscriptions, err := strg.ReadUserSubscriptions(subscription.User{ID: "user-1"})
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_readUserAllSubscription(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)

	_, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  11,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "333-222",
				SizeID:  15,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-3"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	subscriptions, err := strg.ReadUserAllSubscriptions(subscription.User{ID: "user-1"})
	assert.NoError(err)
	assert.Equal(3, len(subscriptions))
}

func TestStorageMemory_readSubscriptionsByShopItem(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)

	_, err := strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  15,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "333-222",
				SizeID:  15,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-3"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	subscriptions, err := strg.ReadSubscriptionsByShopItem(shop.Item{Article: "222-222", SizeID: 10})
	assert.NoError(err)

	assert.Equal(2, len(subscriptions))

	expected := []subscription.Item{
		{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  10,
			},
		},
		{
			Active: true,
			User:   subscription.User{ID: "user-3"},
			ShopItem: shop.Item{
				Article: "222-222",
				SizeID:  10,
			},
		},
	}
	assert.ElementsMatch(expected, subscriptions)
}

func TestStorageMemory_disableSubscription(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-5"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)

	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	err = strg.DisableSubscription(
		subscription.Item{
			User:     subscription.User{ID: "user-1"},
			ShopItem: shop.Item{Article: "111-222", SizeID: 10},
		},
	)
	assert.NoError(err)

	subscriptions, err = strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	subscriptions, err = strg.ReadAllSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMemory_removeSubscription(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	subscriptions, err := strg.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(0, len(subscriptions))

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: true,
			User:   subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)
	assert.NoError(err)

	_, err = strg.CreateSubscription(
		subscription.Item{
			Active: false,
			User:   subscription.User{ID: "user-5"},
			ShopItem: shop.Item{
				Article: "222-123",
				SizeID:  10,
			},
		},
	)

	assert.NoError(err)

	removed, err := strg.RemoveSubscription(
		subscription.Item{
			User:     subscription.User{ID: "user-1"},
			ShopItem: shop.Item{Article: "111-222", SizeID: 10},
		},
	)
	assert.NoError(err)
	assert.True(removed)

	subscriptions, err = strg.ReadAllSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))
}
