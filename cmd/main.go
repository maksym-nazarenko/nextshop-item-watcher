package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/system"

	"github.com/spf13/viper"
)

func loadConfig() (system.Config, error) {
	var config system.Config

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.SetEnvPrefix("NWI")
	viper.AutomaticEnv()

	keysToBind := []string{
		"http.client.baseurl",
		"http.client.lang",
		"watch.updateinterval",
		"bot.allowedusers",
		"bot.token",
		"storage.driver",
		"storage.options",
	}
	if err := func(keys []string) error {
		for _, k := range keys {
			if errViper := viper.BindEnv(k); errViper != nil {
				return errViper
			}
		}
		return nil
	}(keysToBind); err != nil {
		return config, err
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
		// config file not found, but keep trying with Env
		log.Printf("Config file not found: %s", err.Error())
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	// todo(maks): introduce interface Validator.Validate() for configuration
	if config.Storage.Driver == "" {
		return config, errors.New("storage section must contain 'driver' key")
	}

	return config, nil
}

func runSystem(config system.Config) int {
	system := system.New(config)

	if err := system.Start(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	gracefulShutdownInProgress := false

	for {
		select {
		case s := <-c:
			if s == os.Interrupt {
				if gracefulShutdownInProgress {
					log.Fatalln("Forced shutdown requested. Exiting immediately.")
					return 1
				}

				log.Printf("Got signal: %s\n", s)
				log.Println("Performing graceful shutdown. Press Ctrl+C again to force.")
				gracefulShutdownInProgress = true

				system.Stop()
			}

		case <-system.StoppedCh():
			log.Println("[INFO] Exiting")
			return 0
		}
	}
}

func realMain() int {
	ret := 0

	config, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	ret = runSystem(config)

	return ret
}

func main() {
	os.Exit(realMain())
}
