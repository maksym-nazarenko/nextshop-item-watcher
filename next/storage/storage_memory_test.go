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
		subscription.Item{
			User: subscription.User{ID: "user-1"},
			ItemOption: shop.ItemOption{
				Article:           "111-222",
				Name:              "Item 1",
				Price:             "10 euros",
				StockStatusString: shop.ItemStatusInStock,
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
		subscription.Item{
			User: subscription.User{ID: "user-1"},
			ItemOption: shop.ItemOption{
				Article:           "111-222",
				Name:              "Item 1",
				Price:             "10 euros",
				Number:            10,
				StockStatusString: shop.ItemStatusInStock,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))

	added, err = strg.CreateSubscription(
		subscription.Item{
			User: subscription.User{ID: "user-1"},
			ItemOption: shop.ItemOption{
				Article:           "111-222",
				Name:              "Item 1",
				Price:             "10 euros",
				Number:            10,
				StockStatusString: shop.ItemStatusInStock,
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
		subscription.Item{
			User: subscription.User{ID: "user-1"},
			ItemOption: shop.ItemOption{
				Article:           "111-333",
				Name:              "Item 1",
				Price:             "11 euros",
				Number:            11,
				StockStatusString: shop.ItemStatusInStock,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))

	added, err = strg.CreateSubscription(
		subscription.Item{
			User: subscription.User{ID: "user-1"},
			ItemOption: shop.ItemOption{
				Article:           "111-333",
				Name:              "Item 2",
				Price:             "12 euros",
				Number:            12,
				StockStatusString: shop.ItemStatusInStock,
			},
		},
	)

	assert.True(added)
	assert.NoError(err)

	assert.Equal(1, len(strg.ReadSubscriptions()))
}
