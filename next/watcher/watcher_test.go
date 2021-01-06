package watcher

import (
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

func TestAddHandlers_addsNoHandlersIfEmptyVarargsPassed(t *testing.T) {
	watcher := ItemWatcher{}
	assert.Len(t, watcher.handlers, 0)

	watcher.AddHandlers()
	assert.Len(t, watcher.handlers, 0)
}
