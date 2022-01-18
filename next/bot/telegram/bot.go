package telegram

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/shop"

	"github.com/mitchellh/mapstructure"

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
	b.tb.Handle("/help", func(msg *telebot.Message) {
		b.updateBotCommands()
	})
	b.tb.Handle(telebot.OnCallback, b.callbackDispatcher)
	b.tb.Handle(telebot.OnText, b.cmdNewArticle)

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

func (b *Bot) callbackDispatcher(c *telebot.Callback) {
	dataItems := strings.Split(c.Data, "|")
	if len(dataItems) != 2 {
		log.Printf("[ERROR] expected data of %d items, got %d\n", 2, len(dataItems))
		return
	}

	callbackData := NewCallbackData()
	decodedData, err := callbackData.Decode(dataItems[1])
	if err != nil {
		log.Println("[ERROR] Could not decode the callback data: " + err.Error())
		return
	}

	var inlineCallbackData struct {
		Article string
		Size    int
	}

	if err := mapstructure.Decode(decodedData, &inlineCallbackData); err != nil {
		log.Println("[ERROR] Could not map the callback data to structure: " + err.Error())
	}

	created, err := b.mediator.CreateSubscription(
		subscription.Item{
			Active: true,
			User: subscription.User{
				ID: strconv.FormatInt(c.Sender.ID, 10),
			},
			ShopItem: shop.NewItem(inlineCallbackData.Article, inlineCallbackData.Size),
		},
	)

	if err != nil {
		log.Println("[ERROR] subscription creation failed: " + err.Error())
		if _, err = b.tb.Edit(c.Message, "Subscription creation failed"); err != nil {
			log.Println("[ERROR] Could not update message: " + err.Error())
		}

		return
	}

	var messageText string
	if created {
		messageText = fmt.Sprintf("Subscription for %s with sizeID %d created",
			inlineCallbackData.Article,
			inlineCallbackData.Size)
	} else {
		messageText = fmt.Sprintf("Subscription for %s with sizeID %d creation skipped: exists",
			inlineCallbackData.Article,
			inlineCallbackData.Size)
	}

	if _, err = b.tb.Edit(c.Message, messageText); err != nil {
		log.Println("[ERROR] Could not update message: " + err.Error())
	}
}

func (b *Bot) handleInStockItem(item subscription.Item) {
	log.Println("[DEBUG] Bot: new item in stock: ", item)
	inlineURL := &telebot.ReplyMarkup{}

	if item.ShopItem.URL != "" {
		inlineURL.Inline(
			inlineURL.Row(
				inlineURL.URL(item.ShopItem.Description+", "+item.ShopItem.SizeString, item.ShopItem.URL),
			),
		)
	}
	_, err := b.tb.Send(
		ChatID(item.User.ID),
		"Item in stock",
		inlineURL,
	)
	if err != nil {
		log.Println("[ERROR] Could not notify user about in-stock item: " + err.Error())
	}
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

	b.updateBotCommands()
}

func (b *Bot) cmdNewArticle(m *telebot.Message) {
	article, err := ParseStringWithArticle(m.Text)

	if err != nil {
		log.Println("[ERROR] Could not parse article: " + err.Error())
		if _, err := b.tb.Reply(m, err.Error()); err != nil {
			log.Println("[ERROR] Could send message: " + err.Error())
		}

		return
	}

	items, err := b.mediator.FetchSizeIDs(article)
	if err != nil {
		log.Println("[ERROR] Could fetch sized: " + err.Error())
		if _, err := b.tb.Send(m.Sender, "Could not fetch sizes fo "+article); err != nil {
			log.Println("[ERROR] Could send message: " + err.Error())
		}
	}

	inlineSizeSelector := &telebot.ReplyMarkup{}
	rows := make([]telebot.Row, 0, len(items))

	for _, item := range items {
		callbackData := NewCallbackData()
		callbackData.AddItem("article", article)
		callbackData.AddItem("size", item.Number)

		encodedData, err := callbackData.Encode()
		if err != nil {
			log.Println("[ERROR] Could not encode button data: " + err.Error())
		}

		rows = append(
			rows,
			inlineSizeSelector.Row(
				inlineSizeSelector.Data(item.Name, "article="+article+"_size="+strconv.Itoa(item.Number), encodedData),
			))
	}

	inlineSizeSelector.Inline(rows...)

	if _, err := b.tb.Send(m.Sender, "Select size for article "+article, inlineSizeSelector); err != nil {
		log.Println("[ERROR] Could send message: " + err.Error())
		return
	}
}

func (b *Bot) updateBotCommands() {
	log.Println("[INFO] Updating bot commands")
	err := b.tb.SetCommands(
		[]telebot.Command{
			{Text: "/new", Description: "Create new subscription"},
			{Text: "/list", Description: "List all my active subscriptions"},
			{Text: "/help", Description: "Show help"},
		},
	)

	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
	}
}

// New instantiates new Bot object
func New(httpClient *next.Client, mediator *mediator.SubscriptionMediator, config *Config) (*Bot, error) {
	if config.Token == "" {
		return nil, errors.New("telegram Bot token must be set")
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
