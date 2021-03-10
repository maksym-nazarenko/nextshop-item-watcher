package storage

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

type mongoServer struct {
	mongodbC testcontainers.Container
	Port     int
}

func (s *mongoServer) start() error {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017"),
	}

	ctx := context.Background()
	var err error
	s.mongodbC, err = testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return err
	}

	p, err := s.mongodbC.MappedPort(ctx, "27017")
	if err != nil {
		return errors.New("Can't get MongoDB container mapped port")
	}

	s.Port = p.Int()

	return nil
}

func (s *mongoServer) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return s.mongodbC.Terminate(ctx)
}

func createContainer() (*mongoServer, error) {
	tm := mongoServer{}
	if err := tm.start(); err != nil {
		return nil, err
	}

	return &tm, nil
}

func TestStorageMongo_addItemToEmptyStorage(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()

	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_addItemAddsSameItemOnlyOnce(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))

	assert.NoError(err)

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

func TestStorageMongo_addItemAddsSecondItemIfDifferent(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()

	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_readSubscriptionsFetchesOnlyActive(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_enableSubscription(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_readAllSubscriptions(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

	subscriptions, err := strg.ReadAllSubscriptions()
	assert.NoError(err)
	assert.Equal(2, len(subscriptions))
}

func TestStorageMongo_readUserSubscription(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_readUserAllSubscription(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_readSubscriptionsByShopItem(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_disableSubscription(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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

func TestStorageMongo_removeSubscription(t *testing.T) {
	assert := assert.New(t)

	container, err := createContainer()
	assert.NoError(err)
	defer func() {
		if err := container.stop(); err != nil {
			panic("Cannot stop mongo container")
		}
	}()
	strg, err := NewMongo(fmt.Sprintf("mongodb://127.0.0.1:%d", container.Port))
	assert.NoError(err)

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
