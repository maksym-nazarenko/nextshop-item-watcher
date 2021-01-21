package subscription

import "github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

type User struct {
	ID string
}

type Item struct {
	User       User
	ItemOption shop.ItemOption
}
