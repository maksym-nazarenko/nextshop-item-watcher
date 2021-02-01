package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddItem_addsNewItem(t *testing.T) {
	assert := assert.New(t)

	cd := NewCallbackData()

	cd.AddItem("new int value", 1)
	cd.AddItem("new string value", "a string")

	res, err := cd.Encode()
	assert.NoError(err)

	var decoded map[string]interface{}

	decoded, err = cd.Decode(res)
	assert.NoError(err)
	assert.Len(decoded, 2)
	assert.Equal("a string", decoded["new string value"].(string))
}
