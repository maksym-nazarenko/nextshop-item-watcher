package mediator

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

type Observer struct {
	ID      string
	handler func(subscription.Item)
}

func (o *Observer) GetID() string {
	return o.ID
}

func (o *Observer) Update(item subscription.Item) {
	o.handler(item)
}
