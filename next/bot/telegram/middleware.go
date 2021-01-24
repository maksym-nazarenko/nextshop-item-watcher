package telegram

import (
	"fmt"
	"log"

	"github.com/maxim-nazarenko/nextshop-item-watcher/next/subscription"
	"gopkg.in/tucnak/telebot.v2"
)

// WhitelistUserMiddlewarePoller filters user with whitelisting strategy
type WhitelistUserMiddlewarePoller struct {
	allowedUsers map[string]bool
}

func (p *WhitelistUserMiddlewarePoller) filter(u *telebot.Update) bool {
	if u.Message == nil {
		return true
	}

	if _, ok := p.allowedUsers[fmt.Sprint(u.Message.Sender.ID)]; !ok {
		log.Printf("[INFO] User <%v> is not allowed to contact the bot.", u.Message.Sender.ID)
		return false
	}

	return true
}

// NewWhitelistUserMiddlewarePoller instantiates new WhitelistUserMiddlewarePoller object
func NewWhitelistUserMiddlewarePoller(originalPoller telebot.Poller, allowedUsers []subscription.User) *telebot.MiddlewarePoller {
	allowedMap := make(map[string]bool)
	for _, u := range allowedUsers {
		allowedMap[u.ID] = true
	}

	whitelistMiddleware := WhitelistUserMiddlewarePoller{allowedUsers: allowedMap}

	return telebot.NewMiddlewarePoller(originalPoller, whitelistMiddleware.filter)
}
