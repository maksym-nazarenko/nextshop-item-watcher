package mediator

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

// Observer struct wraps a handler func and provides subscription.Observer interface
type Observer struct {
	ID      string
	handler func(subscription.Item)
}

// GetID impelements subscription.Observer interface
func (o *Observer) GetID() string {
	return o.ID
}

// Update implements subscription.Observer interface delegating the call to handler func
func (o *Observer) Update(item subscription.Item) {
	o.handler(item)
}
