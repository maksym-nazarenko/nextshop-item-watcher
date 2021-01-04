package next

import (
	"io"
	"net/http"
	"testing"
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
			},
			{
				"OptionNumber": "13",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU L стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "14",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU XL стандартный",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "17",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU S для высоких",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "18",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU M для высоких",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "19",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU L для высоких",
				"Price": "635 грн",
				"LinkedItem": []
			},
			{
				"OptionNumber": "20",
				"StockStatus": "ComingSoon",
				"StockMessage": "середина января",
				"OptionName": "EU XL для высоких",
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

	if len(options) != 9 {
		t.Error("options' len must be 9")
	}
}
