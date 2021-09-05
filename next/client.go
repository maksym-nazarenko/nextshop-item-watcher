package next

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
)

const (
	// EndpointGetExtendedOptions is an endpoint for extended options retrieval
	EndpointGetExtendedOptions = "/itemstock/getextendedoptions"

	// EndpointSearch is an endpoint to search items
	EndpointSearch = "/search"
)

// Config holds necessary configuration for HTTPClient
type Config struct {
	BaseURL string
	Lang    string
}

// HTTPClient interface to be implemented by different clients
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type NextClient interface {
	GetOptionsByArticle(article string) ([]shop.ItemOption, error)
	GetItemOption(article string, size int) (shop.ItemOption, error)
	GetItemExtendedOption(article string) (shop.ItemExtendedOption, error)
	FindOptionBySize(options []shop.ItemOption, size int) (shop.ItemOption, bool)
	GetItemURLByArticle(article string) (string, error)
}

// Client is a wrapper for some Next APIs
type Client struct {
	HTTPClient HTTPClient
	BaseURL    string
	Language   string
}

func (c *Client) buildEndpointURL(ep string, pathVars ...string) string {
	endpoint := c.BaseURL + "/" + c.Language + ep
	if len(pathVars) > 0 {
		return endpoint + "/" + strings.Join(pathVars, "/")
	}

	return endpoint
}

// GetOptionsByArticle fetched item options by article
func (c *Client) GetOptionsByArticle(article string) ([]shop.ItemOption, error) {
	extendedOptions, err := c.GetItemExtendedOption(article)

	if err != nil {
		return nil, err
	}

	return extendedOptions.Options, nil
}

// GetItemOption fetches a single option object for article and size combination
func (c *Client) GetItemOption(article string, size int) (shop.ItemOption, error) {
	var option shop.ItemOption

	items, err := c.GetOptionsByArticle(article)
	if err != nil {
		return shop.ItemOption{}, err
	}

	item, found := c.FindOptionBySize(items, size)
	if !found {
		return option, errors.New("item not found")
	}

	item.Article = article

	return item, nil
}

// GetItemExtendedOption fetches available options information for particular article
func (c *Client) GetItemExtendedOption(article string) (shop.ItemExtendedOption, error) {
	url := fmt.Sprintf("%s?_=%d", c.buildEndpointURL(EndpointGetExtendedOptions, article), time.Now().Unix())
	var resp *http.Response
	var err error
	if resp, err = c.HTTPClient.Get(url); err != nil {
		return shop.ItemExtendedOption{},
			fmt.Errorf("couldn't get extended options for <%s> article: %s", article, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return shop.ItemExtendedOption{}, errors.New("can't read the response body")
	}
	resp.Body.Close()

	bodyString := string(body)

	var optionResponse shop.ItemExtendedOption
	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&optionResponse); err != nil {
		return shop.ItemExtendedOption{}, err
	}

	return optionResponse, nil
}

func (c *Client) FindOptionBySize(options []shop.ItemOption, size int) (shop.ItemOption, bool) {
	for _, item := range options {
		if item.Number == size {
			return item, true
		}
	}

	return shop.ItemOption{}, false
}

func (c *Client) GetItemURLByArticle(article string) (string, error) {
	url := fmt.Sprintf("%s?w=%s", c.buildEndpointURL(EndpointSearch), url.QueryEscape(article))
	response, err := c.HTTPClient.Get(url)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 && response.Request != nil && response.Request.URL.Fragment == article {
		return response.Request.URL.String(), nil
	}

	return "", errors.New("invalid URL format in response")
}

// NewClient creates a Next client
func NewClient(httpClient HTTPClient, c Config) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{HTTPClient: httpClient, BaseURL: c.BaseURL, Language: c.Lang}
}
