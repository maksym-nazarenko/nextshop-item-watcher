package shop

import (
	"fmt"
	"regexp"
)

// Item describes one particular item by its article and size
type Item struct {
	Article string
	SizeID  int
}

// NormalizeArticle normalizes article to meet canonical representation
func NormalizeArticle(article string) string {
	ret := regexp.MustCompile("[^0-9]").ReplaceAllLiteralString(article, "")

	return ret
}

// ItemExtendedOption holds extended option response from Next API
type ItemExtendedOption struct {
	Description string
	Options     []ItemOption
}

// ItemOption holds option data in ItemExtendedOption Options[] slice
type ItemOption struct {
	Article           string
	Name              string `json:"OptionName"`
	Number            int    `json:"OptionNumber,string"`
	Price             string
	StockStatusString string `json:"StockStatus"`
}

const (
	// ItemStatusInStock holds string representation of in-stock item
	ItemStatusInStock = "InStock"

	// ItemStatusComingSoon holds string representation of coming soon item
	ItemStatusComingSoon = "ComingSoon"

	// ItemStatusUnknown is a placeholder for unknown status
	ItemStatusUnknown = "Unknown"
)

func (item ItemOption) String() string {
	return fmt.Sprintf("[%s] %s, %s", item.StockStatusString, item.Name, item.Price)
}

func NewItem(article string, size int) Item {
	return Item{Article: NormalizeArticle(article), SizeID: size}
}
