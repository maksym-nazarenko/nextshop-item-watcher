package shop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeArticle(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{input: "123-456", expected: "123456"},
		{input: "0_5-23-4", expected: "05234"},
		{input: "", expected: ""},
	}

	for _, test := range tests {
		t.Run("Normalize <"+test.input+">", func(innerT *testing.T) {
			assert.Equal(innerT, test.expected, NormalizeArticle(test.input))
		})
	}
}
