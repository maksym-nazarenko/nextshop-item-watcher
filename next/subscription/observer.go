package subscription

// Observer is an interface to be implemented by clients which are interested in notification about InStock items
type Observer interface {
	GetID() string
	Update(Item)
}
