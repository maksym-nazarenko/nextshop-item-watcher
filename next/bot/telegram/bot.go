package telegram

import (
	"errors"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/mediator"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
)

// Bot works as a frontend to the systems
type Bot struct {
	httpClient *next.Client
	mediator   *mediator.SubscriptionMediator
	config     *Config
	tb         *telebot.Bot
	stopCh     chan struct{}
}

// Start begins the message loop
func (b *Bot) Start() {
	b.tb.Handle("/start", b.cmdStart)

	b.tb.Handle(telebot.OnText, func(msg *telebot.Message) {

	})

	go func() {
		for {
			select {
			case item := <-b.mediator.InStockItemCh():
				b.handleInStockItem(item)
			case <-b.stopCh:
				return
			}
		}
	}()

	b.tb.Start()
}

func (b *Bot) Stop() {
	log.Println("[INFO] Stopping Telegram bot")
	close(b.stopCh)
	b.tb.Stop()
}

func (b *Bot) handleInStockItem(item subscription.Item) {
	log.Println("[DEBUG] Bot: new item in stock: ", item)
}

func (b *Bot) cmdStart(m *telebot.Message) {
	if !m.Private() {
		return
	}

	log.Printf("[Debug] User <%d> started a conversation.\n", m.Chat.ID)
	_, err := b.tb.Send(m.Sender, "Nice to see you here!")

	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return
	}

}

// New instantiates new Bot object
func New(httpClient *next.Client, mediator *mediator.SubscriptionMediator, config *Config) (*Bot, error) {
	if config.Token == "" {
		return nil, errors.New("Telegram Bot token must be set")
	}

	longPoller := &telebot.LongPoller{
		Timeout: 5 * time.Second,
	}

	tb, err := telebot.NewBot(telebot.Settings{
		Token:  config.Token,
		Poller: NewWhitelistUserMiddlewarePoller(longPoller, config.AllowedUsers),
	})

	if err != nil {
		return nil, err
	}

	bot := &Bot{
		httpClient: httpClient,
		mediator:   mediator,
		config:     config,
		tb:         tb,
		stopCh:     make(chan struct{}),
	}

	return bot, nil
}
