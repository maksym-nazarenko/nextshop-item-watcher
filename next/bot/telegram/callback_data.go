package telegram

import (
	"encoding/json"
)

type CallbackData struct {
	items map[string]interface{}
}

func (cd *CallbackData) Decode(s string) (map[string]interface{}, error) {
	ret := make(map[string]interface{})

	if err := json.Unmarshal([]byte(s), &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (cd *CallbackData) Encode() (string, error) {
	bytes, err := json.Marshal(cd.items)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (cd *CallbackData) AddItem(key string, value interface{}) {
	cd.items[key] = value
}

func NewCallbackData() *CallbackData {
	cd := CallbackData{
		items: make(map[string]interface{}),
	}

	return &cd
}
