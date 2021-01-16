package watch

import (
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
)

// ItemOptionHandler is a common interface for ItemOption processors
type ItemOptionHandler interface {
	Handle(...shop.ItemOption)
}

type handleFuncType func(...shop.ItemOption)

// HandleFuncWrapper simply wraps function to match ItemOptionHandler interface
type HandleFuncWrapper struct {
	handleFunc handleFuncType
}

// Handle handles particular ItemOption
func (h HandleFuncWrapper) Handle(itemOption ...shop.ItemOption) {
	h.handleFunc(itemOption...)
}
