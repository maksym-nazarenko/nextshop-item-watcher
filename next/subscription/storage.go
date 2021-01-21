package subscription

type Reader interface {
	ReadSubscriptions() []Item
}

type Writer interface {
	CreateSubscription(Item) (bool, error)
	RemoveSubscription(Item) (bool, error)
}

type Storage interface {
	Reader
	Writer
}
