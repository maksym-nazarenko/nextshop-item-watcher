package next

import (
	"net/http"
	"testing"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/testutils"

	"github.com/stretchr/testify/assert"
)

func Test_buildEndpointURL_noExtraVars(t *testing.T) {
	c := NewClient(
		testutils.NewClientWithPayload(""),
		Config{
			BaseURL: "https://some.host.com",
			Lang:    "ru",
		},
	)

	assert.Equal(t, "https://some.host.com/ru/v1/endpoint/path", c.buildEndpointURL("/v1/endpoint/path"))
}

func Test_buildEndpointURL_withExtraVars(t *testing.T) {
	c := NewClient(
		testutils.NewClientWithPayload(""),
		Config{
			BaseURL: "https://some.host.com",
			Lang:    "ru",
		},
	)

	c.buildEndpointURL("/v1/endpoint/path", "var1", "var2")

	assert.Equal(t, "https://some.host.com/ru/v1/endpoint/path/var1/var2", c.buildEndpointURL("/v1/endpoint/path", "var1", "var2"))
}

func TestGetOptionsByArticle(t *testing.T) {
	payload := `
	{
		"Description": "Розовая в цветочек - Теплая пижама",
		"ItemNumber": "821-585",
		"ComingSoonEnabled": true,
		"Options": [
			{
				"OptionNumber": "10",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU XS стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "11",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU S стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "12",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU M стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			}
		],
		"PersonalisedGift": "N",
		"PersonalisedGiftTheme": "0",
		"DDFulfiller": "",
		"FulfilmentType": ""
	}`
	client := NewClient(
		testutils.NewClientWithPayload(payload),
		Config{
			BaseURL: "https://www.next.ua",
			Lang:    "ru",
		},
	)

	options, err := client.GetOptionsByArticle("127001")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 3, len(options))
}

func TestGetItemInfo(t *testing.T) {
	payload := `
	{
		"Description": "Розовая в цветочек - Теплая пижама",
		"ItemNumber": "821-585",
		"ComingSoonEnabled": true,
		"Options": [
			{
				"OptionNumber": "10",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU XS стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "11",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU S стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "12",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU M стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			}
		],
		"PersonalisedGift": "N",
		"PersonalisedGiftTheme": "0",
		"DDFulfiller": "",
		"FulfilmentType": ""
	}`
	client := NewClient(
		testutils.NewClientWithPayload(payload),
		Config{
			BaseURL: "https://www.next.ua",
			Lang:    "ru",
		},
	)

	option, err := client.GetItemOption("821-585", 11)
	assert := assert.New(t)

	assert.NoError(err)
	assert.NotNil(option)
	assert.EqualValues(shop.ItemOption{Article: "821-585", Name: "EU S стандартный", Number: 11, Price: "635 грн", StockStatusString: "ComingSoon"}, option)
}

func TestGetItemInfo_returnsErrorOnWrongSizeID(t *testing.T) {
	payload := `
	{
		"Description": "Розовая в цветочек - Теплая пижама",
		"ItemNumber": "821-585",
		"ComingSoonEnabled": true,
		"Options": [
			{
				"OptionNumber": "10",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU XS стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "11",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU S стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "12",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU M стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			}
		],
		"PersonalisedGift": "N",
		"PersonalisedGiftTheme": "0",
		"DDFulfiller": "",
		"FulfilmentType": ""
	}`
	client := NewClient(
		testutils.NewClientWithPayload(payload),
		Config{
			BaseURL: "https://www.next.ua",
			Lang:    "ru",
		},
	)

	_, err := client.GetItemOption("821-585", 1)
	assert.Error(t, err)
}

func TestNewClient_useDefaultClientIfNotOverridden(t *testing.T) {
	client := NewClient(nil,
		Config{
			BaseURL: "https://www.next.ua",
			Lang:    "ru",
		},
	)
	assert.IsType(t, http.DefaultClient, client.HTTPClient)
}

func TestNewClient_useProvidedClient(t *testing.T) {
	mockedHTTPClient := testutils.NewClientWithPayload("")
	client := NewClient(
		mockedHTTPClient,
		Config{
			BaseURL: "https://www.next.ua",
			Lang:    "ru",
		},
	)
	assert.IsType(t, mockedHTTPClient, client.HTTPClient)
}
