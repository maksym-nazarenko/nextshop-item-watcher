package watch

import (
	"testing"
	"time"

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

	w := New(
		next.NewClient(
			testutils.NewClientWithPayload(payload),
			"https://www.example.com", "ru",
		),
		&Config{UpdateInterval: 1 * time.Second},
	)

	w.AddItem(next.ShopItem{Article: "821-585", SizeID: 11})

	w.Process()
	defer w.Stop()

	select {
	case <-w.InStockChan():
		return
	case <-time.After(time.Second):
		t.Error("No items received")
	}
}
