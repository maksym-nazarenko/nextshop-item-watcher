package mediator

import (
	"testing"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/testutils"

	"github.com/stretchr/testify/assert"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/storage"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

func TestCreateSubscription(t *testing.T) {
	storage := storage.NewMemoryStorage()
	watcher, _ := watch.New(nil, &watch.Config{UpdateInterval: 2 * time.Second})
	mediator := New(
		storage,
		watcher,
		next.NewClient(
			testutils.NewClientWithPayload(""),
			next.Config{
				BaseURL: "",
				Lang:    "uk",
			},
		),
	)

	assert := assert.New(t)

	item := subscription.Item{
		Active:   true,
		User:     subscription.User{ID: "user-"},
		ShopItem: shop.NewItem("111-222", 10),
	}

	ok, err := mediator.CreateSubscription(
		item,
	)

	assert.NoError(err)
	assert.True(ok)

	subscriptions, err := mediator.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(subscriptions))

	storageSubscriptions, err := storage.ReadSubscriptions()
	assert.NoError(err)
	assert.Equal(1, len(storageSubscriptions))
}
