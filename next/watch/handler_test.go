package watch

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
)

func TestHandleFuncWrapper_wrapsPlainFunction(t *testing.T) {
	var actualString string

	option := next.ItemOption{Name: "test item", Number: 29, Price: "12 euros", StockStatusString: "InStock"}
	expectedString := option.String()

	wrapper := HandleFuncWrapper{handleFunc: func(options ...next.ItemOption) {
		actualString = options[0].String()
	}}

	wrapper.Handle(option)

	assert.Equal(t, expectedString, actualString)
}
