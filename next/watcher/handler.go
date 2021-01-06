package watcher

import "github.com/maxim-nazarenko/nextshop-item-watcher/next"

// ItemOptionHandler is a common interface for ItemOption processors
type ItemOptionHandler interface {
	Handle(...next.ItemOption)
}

type handleFuncType func(...next.ItemOption)

// HandleFuncWrapper simply wraps function to match ItemOptionHandler interface
type HandleFuncWrapper struct {
	handleFunc handleFuncType
}

// Handle handles particular ItemOption
func (h HandleFuncWrapper) Handle(itemOption ...next.ItemOption) {
	h.handleFunc(itemOption...)
}
