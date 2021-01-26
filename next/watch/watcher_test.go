package watch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/testutils"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
)

func TestWatcherPassesInStockItemsToChannel(t *testing.T) {
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
				"StockStatus": "InStock",
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

	w, err := New(
		next.NewClient(
			testutils.NewClientWithPayload(payload),
			next.Config{
				BaseURL: "https://www.next.ua",
				Lang:    "ru",
			},
		),
		&Config{UpdateInterval: 10 * time.Millisecond},
	)

	assert.NoError(t, err)

	err = w.AddItem(&shop.Item{Article: "821-585", SizeID: 11})
	assert.NoError(t, err)

	w.Run()
	defer func() {
		w.Stop()
	}()

	select {
	case <-w.InStockChan():
		return
	case <-time.After(2 * time.Second):
		t.Error("No items received")
	}
}
