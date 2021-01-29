package telegram

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
)

func parseArticleByLink(l string) (string, error) {
	parsedURL, err := url.Parse(l)

	if err != nil {
		return "", err
	}

	if !parsedURL.IsAbs() {
		return "", errors.New("Only absolute URLs are supported")
	}

	if parsedURL.Scheme != "https" {
		return "", errors.New("Only https:// scheme is supported")
	}

	pathComponents := strings.Split(parsedURL.Path, "/")
	article := shop.NormalizeArticle(pathComponents[len(pathComponents)-1])
	if len(article) == 6 {
		return article, nil
	}

	return "", errors.New("Cannot extract article from the link")
}

// ParseStringWithArticle parses a string and extracts the article number
//
// Current implementation supports the following formats:
//	1. Raw article number: 111222, 111-222, 111_222
//	2. Link to shop item page: https://www.domain.com/.../111222
func ParseStringWithArticle(msg string) (string, error) {
	if strings.HasPrefix(msg, "http") {
		return parseArticleByLink(msg)
	}

	article := shop.NormalizeArticle(msg)
	// TODO: move parsing+validation logic to separate data type Article
	if len(article) == 6 {
		return article, nil
	}

	return "", fmt.Errorf("Unknown article format <%s>. Use article number '111-222', '111222' or link to item page", msg)
}
