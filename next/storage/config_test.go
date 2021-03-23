package storage_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/system"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	var configString = []byte(`
http:
  client:
    baseURL: "https://www.next.ua"
    lang: "ru"

watch:
  updateInterval: "3s"

bot:
  allowedUsers:
  - id: "219836184"
token: "secret token"

storage:
  driver: mongo
  options:
    url: "url"
    timeout: "3s"
`)
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(configString))
	assert.NoError(t, err)

	var config system.Config
	err = viper.Unmarshal(&config)
	assert.NoError(t, err)

	type MongoConfig struct {
		Url     string
		Timeout time.Duration
	}

	var mongoConfig MongoConfig
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result:   &mongoConfig,
			Metadata: nil,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
			),
		},
	)
	assert.NoError(t, err)

	err = decoder.Decode(config.Storage.Options)
	assert.NoError(t, err)

	t.Log(config)
	t.Log(mongoConfig)
}
