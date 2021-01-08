package watch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
)

type TestHandler int

func (t TestHandler) Handle(items ...next.ItemOption) {

}

func TestAddHandlers_addsAllHandlers(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlers(new(TestHandler), new(TestHandler))
	assert.Len(t, watcher.handlers, 2)
}

func TestAddHandlers_appendsHandlers(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlers(new(TestHandler), new(TestHandler))
	assert.Len(t, watcher.handlers, 2)

	watcher.AddHandlers(new(TestHandler))
	assert.Len(t, watcher.handlers, 3)
}

func TestAddHandlers_addsNoHandlersIfEmptyVarargsPassed(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlers()
	assert.Len(t, watcher.handlers, 0)
}

func TestAddHandlerFuncs_addsAllHandlerFuncs(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlerFuncs(func(...next.ItemOption) {}, func(...next.ItemOption) {}, func(...next.ItemOption) {})

	assert.Len(t, watcher.handlers, 3)
}

func TestAddHandlerFuncs_appendsHandlerFuncs(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlerFuncs(func(...next.ItemOption) {}, func(...next.ItemOption) {}, func(...next.ItemOption) {})
	assert.Len(t, watcher.handlers, 3)

	watcher.AddHandlerFuncs(func(...next.ItemOption) {}, func(...next.ItemOption) {})
	assert.Len(t, watcher.handlers, 5)
}

func TestProcessInStockItems_callsAllHandlers(t *testing.T) {
	watcher := ItemWatcher{}

	handledItems := []string{}

	watcher.AddHandlerFuncs(func(items ...next.ItemOption) {
		handledItems = append(handledItems, fmt.Sprintf("handler1 handled %d items", len(items)))
	})
}
