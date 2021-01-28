package mediator

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/storage"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
)

func TestCreateSubscription(t *testing.T) {
	storage := storage.NewMemoryStorage()
	watcher, _ := watch.New(nil, &watch.Config{UpdateInterval: 2 * time.Second})
	mediator := New(storage, watcher)

	assert := assert.New(t)

	item := subscription.Item{
		User:     subscription.User{ID: "user-"},
		ShopItem: shop.NewItem("111-222", 10),
	}

	ok, err := item.RegisterObserver(
		&Observer{
			ID: t.Name() + "-observer",
			handler: func(item subscription.Item) {
				log.Printf("[DEBUG] Item is in stock: %v\n", item)
			},
		},
	)

	assert.NoError(err)
	assert.True(ok)

	ok, err = mediator.CreateSubscription(
		&item,
	)

	assert.NoError(err)
	assert.True(ok)

	assert.Equal(1, len(mediator.ReadSubscriptions()))
	assert.Equal(1, len(storage.ReadSubscriptions()))

	assert.Equal(2, len(mediator.ReadSubscriptions()[0].Observers()))
	assert.Equal(2, len(storage.ReadSubscriptions()[0].Observers()))
}
