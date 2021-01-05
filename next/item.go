package next

import "regexp"

// ShopItem describes one particular item by its article and size
type ShopItem struct {
	Article string
	SizeID  int
}

// NormalizeArticle normalizes article to meet canonical representation
func NormalizeArticle(article string) string {
	ret := regexp.MustCompile("[^0-9]").ReplaceAllLiteralString(article, "")

	return ret
}
