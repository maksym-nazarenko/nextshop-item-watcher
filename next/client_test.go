package next

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockBody struct {
	Payload     string
	payloadLeft []byte
}

func NewMockBody(payload string) *MockBody {
	body := MockBody{Payload: payload}
	body.ResetReader()
	return &body
}

func (b *MockBody) Read(p []byte) (n int, err error) {
	if len(b.payloadLeft) < len(p) {
		return copy(p, b.payloadLeft), io.EOF
	}

	n = copy(p, b.payloadLeft[:len(p)])
	b.payloadLeft = b.payloadLeft[n:]

	return n, nil
}

func (b *MockBody) ResetReader() {
	b.payloadLeft = []byte(b.Payload)
}

func (b *MockBody) Close() error {
	return nil
}

type MockHTTPClient struct {
	Body io.ReadCloser
}

func (c *MockHTTPClient) Get(url string) (resp *http.Response, err error) {

	return &http.Response{StatusCode: 200, Status: "OK", Body: c.Body}, nil
}

func Test_buildEndpointURL_noExtraVars(t *testing.T) {
	c := NewClient(&MockHTTPClient{Body: NewMockBody("")}, "https://some.host.com", "ru")

	assert.Equal(t, "https://some.host.com/ru/v1/endpoint/path", c.buildEndpointURL("/v1/endpoint/path"))
}

func Test_buildEndpointURL_withExtraVars(t *testing.T) {
	c := NewClient(&MockHTTPClient{Body: NewMockBody("")}, "https://some.host.com", "ru")

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
	client := NewClient(&MockHTTPClient{Body: NewMockBody(payload)}, "https://www.example.com", "ru")

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
	client := NewClient(&MockHTTPClient{Body: NewMockBody(payload)}, "https://www.example.com", "ru")

	option, err := client.GetItemInfo(ShopItem{Article: "821-585", SizeID: 11})
	assert := assert.New(t)

	assert.NoError(err)
	assert.NotNil(option)
	assert.EqualValues(ItemOption{Name: "EU S стандартный", Number: 11, Price: "635 грн", StockStatusString: "ComingSoon"}, option)
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
	client := NewClient(&MockHTTPClient{Body: NewMockBody(payload)}, "https://www.example.com", "ru")

	_, err := client.GetItemInfo(ShopItem{Article: "821-585", SizeID: 1})
	assert.Error(t, err)
}

func TestNewClient_useDefaultClientIfNotOverridden(t *testing.T) {
	client := NewClient(nil, "https://www.example.com", "ru")
	assert.IsType(t, http.DefaultClient, client.HTTPClient)
}

func TestNewClient_useProvidedClient(t *testing.T) {
	mockedHTTPClient := &MockHTTPClient{Body: NewMockBody("")}
	client := NewClient(mockedHTTPClient, "https://www.example.com", "ru")
	assert.IsType(t, mockedHTTPClient, client.HTTPClient)
}
