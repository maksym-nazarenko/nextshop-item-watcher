package subscription

import "github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

// Item represents a subscription item
type Item struct {
	User     User
	ShopItem shop.Item
}
