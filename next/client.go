package next

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"
)

const (
	// EndpointGetExtendedOptions is an endpoint for extended options retrieval
	EndpointGetExtendedOptions = "/itemstock/getextendedoptions"
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
	url := fmt.Sprintf("%s?_=%d", c.buildEndpointURL(EndpointGetExtendedOptions, article), time.Now().Unix())
	var resp *http.Response
	var err error
	if resp, err = c.HTTPClient.Get(url); err != nil {
		return nil, fmt.Errorf("Couldn't get extended options for <%s> article: %s", article, err.Error())
	}

	var optionResponse shop.ItemExtendedOption

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Can't read the response body")
	}
	resp.Body.Close()

	bodyString := string(body)

	if err := json.NewDecoder(strings.NewReader(bodyString)).Decode(&optionResponse); err != nil {
		return nil, err
	}

	return optionResponse.Options, nil
}

// GetItemInfo checks the state of a particular shop item
func (c *Client) GetItemInfo(shopItem shop.Item) (shop.ItemOption, error) {
	var option shop.ItemOption

	items, err := c.GetOptionsByArticle(shopItem.Article)
	if err != nil {
		return shop.ItemOption{}, err
	}

	for _, item := range items {
		if item.Number == shopItem.SizeID {
			item.Article = shopItem.Article

			return item, nil
		}
	}

	return option, errors.New("Item not found")
}

// NewClient creates a Next client
func NewClient(httpClient HTTPClient, c Config) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{HTTPClient: httpClient, BaseURL: c.BaseURL, Language: c.Lang}
}
