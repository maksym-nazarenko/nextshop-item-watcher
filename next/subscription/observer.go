package subscription

// Observer interface provides
type Observer interface {
	GetID() string
	Update(Item)
}
