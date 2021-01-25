package subscription

import "github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

// Item represents a subscription item
type Item struct {
	User      User
	ShopItem  shop.Item
	observers []Observer
}

// RegisterObserver register new observer for this item
// returns false if the observer was not added,
// e.g it is already registered
func (i *Item) RegisterObserver(o Observer) (bool, error) {
	if i.observerRegistered(o) {
		return false, nil
	}

	i.observers = append(i.observers, o)

	return true, nil
}

// DeregisterObserver deregister observer for this item
// returns false if the observer was previously registered
func (i *Item) DeregisterObserver(o Observer) (bool, error) {
	if !i.observerRegistered(o) {
		return false, nil
	}

	if err := i.doDeregisterObserver(o); err != nil {
		return false, err
	}

	return true, nil
}

// NotifyAll notifies all registered observers
func (i *Item) NotifyAll() {
	for _, o := range i.observers {
		o.Update(*i)
	}
}

// Observers returns a slice of registered observers
func (i *Item) Observers() []Observer {
	return i.observers
}

func (i *Item) observerRegistered(o Observer) bool {
	for _, item := range i.observers {
		if item.GetID() == o.GetID() {
			return true
		}
	}
	return false
}

func (i *Item) doDeregisterObserver(o Observer) error {
	var index int
	for idx, item := range i.observers {
		if item.GetID() == o.GetID() {
			index = idx
			break
		}
	}

	i.observers[index] = i.observers[len(i.observers)-1]
	i.observers = i.observers[:len(i.observers)-1]

	return nil
}
