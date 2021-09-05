package testutils

import "github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

type MockNextClientHandlers struct {
	GetOptionsByArticle   func(article string) ([]shop.ItemOption, error)
	GetItemOption         func(article string, size int) (shop.ItemOption, error)
	GetItemExtendedOption func(article string) (shop.ItemExtendedOption, error)
	FindOptionBySize      func(options []shop.ItemOption, size int) (shop.ItemOption, bool)
	GetItemURLByArticle   func(article string) (string, error)
}
type MockNextClient struct {
	Handlers MockNextClientHandlers
}

func (c *MockNextClient) GetOptionsByArticle(article string) ([]shop.ItemOption, error) {
	return c.Handlers.GetOptionsByArticle(article)
}

func (c *MockNextClient) GetItemOption(article string, size int) (shop.ItemOption, error) {
	return c.Handlers.GetItemOption(article, size)
}

func (c *MockNextClient) GetItemExtendedOption(article string) (shop.ItemExtendedOption, error) {
	return c.Handlers.GetItemExtendedOption(article)
}

func (c *MockNextClient) FindOptionBySize(options []shop.ItemOption, size int) (shop.ItemOption, bool) {
	return c.Handlers.FindOptionBySize(options, size)
}

func (c *MockNextClient) GetItemURLByArticle(article string) (string, error) {
	return c.Handlers.GetItemURLByArticle(article)
}

func NewMockNextClient(handlers MockNextClientHandlers) *MockNextClient {
	return &MockNextClient{Handlers: handlers}
}
