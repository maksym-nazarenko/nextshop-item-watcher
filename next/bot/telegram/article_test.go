package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArticle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Raw article",
			input:    "123456",
			expected: "123456",
		},
		{
			name:     "Raw article with dash",
			input:    "091-291",
			expected: "091291",
		},
		{
			name:     "Link to item page",
			input:    "https://www.domain.com/uk/collections/spring/134857",
			expected: "134857",
		},
		{
			name:     "Link to item page with fragment",
			input:    "https://www.domain.com/uk/collections/spring/234857#anchor",
			expected: "234857",
		},
		{
			name:     "Link to item page with fragment and query",
			input:    "https://www.domain.com/uk/collections/spring/934813?sort=asc#anchor",
			expected: "934813",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			article, err := ParseStringWithArticle(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, article)
		})
	}
}

func TestParseArticle_returnErrorOnUnsupportedFormats(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Raw article with less than 6 symbols",
			input: "12345",
		},
		{
			name:  "Empty article string",
			input: "",
		},
		{
			name:  "Raw article with more than 6 symbols",
			input: "12345678",
		},
		{
			name:  "Link to item page with http scheme",
			input: "http://www.domain.com/uk/collections/spring/234857#anchor",
		},
		{
			name:  "Link to item page with invalid article format",
			input: "https://www.domain.com/uk/collections/spring/93_813",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			article, err := ParseStringWithArticle(test.input)
			assert.Error(t, err)
			assert.Equal(t, "", article)
		})
	}
}
