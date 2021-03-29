module github.com/maxim-nazarenko/nextshop-item-watcher

go 1.14

require (
	github.com/mitchellh/mapstructure v1.1.2
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.10.0
	go.mongodb.org/mongo-driver v1.4.6
	gopkg.in/tucnak/telebot.v2 v2.3.5
)

replace github.com/testcontainers/testcontainers-go v0.10.0 => github.com/maxim-nazarenko/testcontainers-go v0.10.1-0.20210323214920-5d30355e35ce
