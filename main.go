package main

import (
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.SetEnvPrefix("NWI")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
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
