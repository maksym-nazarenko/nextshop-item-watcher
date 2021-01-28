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
	assert.Equal(0, len(strg.ReadSubscriptions()))

	added, err := strg.CreateSubscription(
		&subscription.Item{
			User: subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  1,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))
}

func TestStorageMemory_addItemAddsSameItemOnlyOnce(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	assert.Equal(0, len(strg.ReadSubscriptions()))

	added, err := strg.CreateSubscription(
		&subscription.Item{
			User: subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))

	added, err = strg.CreateSubscription(
		&subscription.Item{
			User: subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-222",
				SizeID:  10,
			},
		},
	)

	assert.False(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))
}

func TestStorageMemory_addItemAddsSecondItemIfDifferent(t *testing.T) {
	strg := NewMemoryStorage()

	assert := assert.New(t)
	assert.Equal(0, len(strg.ReadSubscriptions()))

	added, err := strg.CreateSubscription(
		&subscription.Item{
			User: subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-333",
				SizeID:  11,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))

	added, err = strg.CreateSubscription(
		&subscription.Item{
			User: subscription.User{ID: "user-1"},
			ShopItem: shop.Item{
				Article: "111-333",
				SizeID:  12,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(2, len(strg.ReadSubscriptions()))
}
