package testutils

import (
	"io"
	"net/http"
)

type MockBody struct {
	Payload     string
	payloadLeft []byte
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

func NewMockBody(payload string) *MockBody {
	body := MockBody{Payload: payload}
	body.ResetReader()
	return &body
}

func NewClientWithPayload(payload string) *MockHTTPClient {
	return &MockHTTPClient{Body: NewMockBody(payload)}
}
