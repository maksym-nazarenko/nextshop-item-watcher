package storage

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"

	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionUser struct {
	ID string
}
type ShopItem struct {
	Article     string
	SizeID      int
	Description string
	SizeString  string
	URL         string
}
type SubscriptionItem struct {
	Active   bool
	ShopItem ShopItem
	User     SubscriptionUser
}

type MongoStorage struct {
	client *mongo.Client
}

func (m *MongoStorage) ReadSubscriptions() ([]subscription.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.client.Database("next").Collection("subscriptions").Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var returnItems []subscription.Item
	err = cursor.All(ctx, &returnItems)
	if err != nil {
		return nil, err
	}

	return returnItems, nil
}

func (m *MongoStorage) ReadAllSubscriptions() ([]subscription.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.client.Database("next").Collection("subscriptions").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var returnItems []subscription.Item
	err = cursor.All(ctx, &returnItems)
	if err != nil {
		return nil, err
	}

	return returnItems, nil
}

func (m *MongoStorage) ReadUserSubscriptions(user subscription.User) ([]subscription.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.client.Database("next").Collection("subscriptions").Find(ctx, bson.M{"user.id": user.ID, "active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var returnItems []subscription.Item
	err = cursor.All(ctx, &returnItems)
	if err != nil {
		return nil, err
	}

	return returnItems, nil
}

func (m *MongoStorage) ReadUserAllSubscriptions(user subscription.User) ([]subscription.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.client.Database("next").Collection("subscriptions").Find(ctx, bson.M{"user.id": user.ID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var returnItems []subscription.Item
	err = cursor.All(ctx, &returnItems)
	if err != nil {
		return nil, err
	}

	return returnItems, nil
}

func (m *MongoStorage) ReadSubscriptionsByShopItem(item shop.Item) ([]subscription.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.client.Database("next").Collection("subscriptions").Find(ctx, bson.M{"shopitem.article": item.Article, "shopitem.sizeid": item.SizeID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var returnItems []subscription.Item
	err = cursor.All(ctx, &returnItems)
	if err != nil {
		return nil, err
	}

	return returnItems, nil
}

func (m *MongoStorage) CreateSubscription(item subscription.Item) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := m.client.Database("next").Collection("subscriptions").FindOne(
		ctx,
		bson.D{
			{"user.id", item.User.ID},
			{"shopitem.article", item.ShopItem.Article},
			{"shopitem.sizeid", item.ShopItem.SizeID},
		},
	)

	if res.Err() == nil {
		return false, nil
	}

	if !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return false, res.Err()
	}

	_, err := m.client.Database("next").Collection("subscriptions").InsertOne(ctx, &item)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *MongoStorage) DisableSubscription(item subscription.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.Database("next").Collection("subscriptions").UpdateOne(
		ctx,
		bson.M{"shopitem.article": item.ShopItem.Article, "shopitem.sizeid": item.ShopItem.SizeID, "user.id": item.User.ID},
		bson.D{
			{"$set", bson.M{"active": false}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStorage) EnableSubscription(item subscription.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.Database("next").Collection("subscriptions").UpdateOne(
		ctx,
		bson.M{"shopitem.article": item.ShopItem.Article, "shopitem.sizeid": item.ShopItem.SizeID, "user.id": item.User.ID},
		bson.D{
			{"$set", bson.M{"active": true}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStorage) RemoveSubscription(item subscription.Item) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.Database("next").Collection("subscriptions").DeleteOne(
		ctx,
		bson.M{"shopitem.article": item.ShopItem.Article, "shopitem.sizeid": item.ShopItem.SizeID, "user.id": item.User.ID},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *MongoStorage) disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func NewMongo(uri string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	mongoStorage := MongoStorage{client: client}

	return &mongoStorage, nil
}
