package subscription

// Reader describes interface for read operations
type Reader interface {
	ReadSubscriptions() []Item
}

// Writer describes interface for write operations
type Writer interface {
	CreateSubscription(Item) (bool, error)
	RemoveSubscription(Item) (bool, error)
}

// Storage is an interface for all storage-related operations on subscriptions
type Storage interface {
	Reader
	Writer
}
