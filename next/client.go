package next

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

// ItemExtendedOption holds extended option response from Next API
type ItemExtendedOption struct {
	Description string
	Options     []ItemOption
}

// ItemOption holds option data in ItemExtendedOption Options[] slice
type ItemOption struct {
	Name              string `json:"OptionName"`
	Number            int    `json:"OptionNumber,string"`
	Price             string
	StockStatusString string `json:"StockStatus"`
	Status            *StockStatus
}

const (
	// ItemStatusInStock holds string representation of in-stock item
	ItemStatusInStock = "InStock"

	// ItemStatusComingSoon holds string representation of coming soon item
	ItemStatusComingSoon = "ComingSoon"

	// ItemStatusUnknown is a placeholder for unknown status
	ItemStatusUnknown = "Unknown"
)

const (
	// EndpointGetExtendedOptions is an endpoint for extended options retrieval
	EndpointGetExtendedOptions = "/itemstock/getextendedoptions"
)

// StockStatus holds various statuses like "InStock", "ComingSoon", etc
type StockStatus string

// Parse parses string representation of status and returns StockStatus type
// or ItemStatusUnknown if status string is invalid
func (s *StockStatus) Parse(value string) *StockStatus {
	var ret StockStatus = ItemStatusUnknown
	switch value {
	case ItemStatusInStock:
		ret = ItemStatusInStock
	case ItemStatusComingSoon:
		ret = ItemStatusComingSoon
	}

	return &ret
}

func (c *Client) buildEndpointURL(ep string, pathVars ...string) string {

	endpoint := c.BaseURL + "/" + c.Language + ep
	if len(pathVars) > 0 {
		return endpoint + "/" + strings.Join(pathVars, "/")
	}

	return endpoint
}

// GetOptionsByArticle fetched item options by article
func (c *Client) GetOptionsByArticle(article string) ([]ItemOption, APIError) {
	url := fmt.Sprintf("%s?_=%d", c.buildEndpointURL(EndpointGetExtendedOptions, article), time.Now().Unix())
	var resp *http.Response
	var err error
	if resp, err = c.HTTPClient.Get(url); err != nil {
		return nil, fmt.Errorf("Couldn't get extended options for <%s> article: %s", article, err.Error())
	}

	var optionResponse ItemExtendedOption

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Can't read the response body")
	}
	resp.Body.Close()

	bodyString := string(body)

	json.NewDecoder(strings.NewReader(bodyString)).Decode(&optionResponse)

	return optionResponse.Options, nil
}

// NewClient creates a Next client
func NewClient(httpClient HTTPClient, baseURL string, lang string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{HTTPClient: httpClient, BaseURL: baseURL, Language: lang}
}
