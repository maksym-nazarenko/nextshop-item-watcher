package telegram

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next"
	"github.com/maxim-nazarenko/nextshop-item-watcher/next/mediator"
)

// Bot works as a frontend to the systems
type Bot struct {
	httpClient *next.Client
	mediator   *mediator.SubscriptionMediator
	config     *Config
	tb         *telebot.Bot
}

// Start begins the message loop
func (b *Bot) Start() {
	b.tb.Handle("/start", b.cmdStart)

	b.tb.Handle(telebot.OnText, func(msg *telebot.Message) {

	})

	b.tb.Start()
}

func (b *Bot) Stop() {
	log.Println("[INFO] Stopping Telegram bot")
	b.tb.Stop()
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
func New(httpClient *next.Client, mediator *mediator.SubscriptionMediator, config *Config, token string) (*Bot, error) {
	longPoller := &telebot.LongPoller{
		Timeout: 5 * time.Second,
	}

	tb, err := telebot.NewBot(telebot.Settings{
		Token:  token,
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
	}

	return bot, nil
}
