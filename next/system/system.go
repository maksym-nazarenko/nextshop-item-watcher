package system

import (
	"log"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/bot/telegram"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/mediator"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/storage"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/watch"
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

func (s *System) doStart() error {
	httpClient := newHTTPCLient(s.config)

	watcher, err := newWatcher(httpClient, s.config)
	if err != nil {
		return err
	}

	if err := watcher.Start(); err != nil {
		return err
	}

	// TODO: dynamically construct the storage
	storage := storage.NewMemoryStorage()
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
