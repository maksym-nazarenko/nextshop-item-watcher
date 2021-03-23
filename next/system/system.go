package system

import (
	"fmt"
	"log"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/bot/telegram"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/mediator"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/storage"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
	"github.com/mitchellh/mapstructure"
)

type StopEvent struct{}

type System struct {
	config    Config
	stopCh    chan struct{}
	stoppedCh chan StopEvent
}

func (s *System) StoppedCh() <-chan StopEvent {
	return s.stoppedCh
}

// Start initializes all the components and runs the processing loop
func (s *System) Start() error {
	s.init()

	return s.doStart()
}

// Stop stops and uninitialized system components
func (s *System) Stop() {
	s.doStop()
}

func (s *System) init() {
	s.stopCh = make(chan struct{}, 1)
	s.stoppedCh = make(chan StopEvent, 1)
}

func (s *System) doStop() {
	close(s.stopCh)
}

func newHTTPCLient(c Config) *next.Client {
	httpClient := next.NewClient(
		nil,
		c.HTTP.Client,
	)

	return httpClient
}

func newWatcher(httpClient *next.Client, c Config) (*watch.ItemWatcher, error) {
	w, err := watch.New(httpClient, &c.Watch)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func newTelegramBot(httpClient *next.Client, mediator *mediator.SubscriptionMediator, c Config) (*telegram.Bot, error) {
	bot, err := telegram.New(
		httpClient,
		mediator,
		&c.Bot,
	)

	if err != nil {
		return nil, err
	}

	return bot, nil
}

func newStorage(storageConfig storage.Config) (storage.Storage, error) {
	supportedDrivers := map[string]func() (storage.Storage, error){
		"memory": func() (storage.Storage, error) {
			return storage.NewMemoryStorage(), nil
		},
		"mongo": func() (storage.Storage, error) {
			var mongoConfig storage.MongoConfig
			decoder, err := mapstructure.NewDecoder(
				&mapstructure.DecoderConfig{
					Result:   &mongoConfig,
					Metadata: nil,
					DecodeHook: mapstructure.ComposeDecodeHookFunc(
						mapstructure.StringToTimeDurationHookFunc(),
					),
				},
			)

			if err != nil {
				return nil, err
			}

			if err := decoder.Decode(storageConfig.Options); err != nil {
				return nil, err
			}

			return storage.NewMongoWithConfig(mongoConfig)
		},
	}

	factoryFunc, ok := supportedDrivers[storageConfig.Driver]
	if !ok {
		return nil, fmt.Errorf("storage driver '%s' is not supported", storageConfig.Driver)
	}
	log.Printf("Initializing '%s' storage", storageConfig.Driver)

	return factoryFunc()
}

func (s *System) doStart() error {
	httpClient := newHTTPCLient(s.config)

	watcher, err := newWatcher(httpClient, s.config)
	if err != nil {
		return err
	}

	if err := watcher.Start(); err != nil {
		return err
	}

	storage, err := newStorage(s.config.Storage)
	if err != nil {
		return err
	}

	mediator := mediator.New(storage, watcher, httpClient)

	bot, err := newTelegramBot(httpClient, mediator, s.config)
	if err != nil {
		return err
	}

	go bot.Start()

	go mediator.Start()

	go func() {
		<-s.stopCh
		defer close(s.stoppedCh)

		log.Println("[INFO] Waiting for all subsystems to shut down")
		mediator.Stop()
		bot.Stop()
		watcher.Stop()
		log.Println("[INFO] All subsystems are shut down")
		s.stoppedCh <- StopEvent{}
	}()

	return nil
}

func New(config Config) *System {
	return &System{config: config}
}
