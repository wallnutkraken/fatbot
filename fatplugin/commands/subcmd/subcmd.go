package subcmd

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/fatbot/fatbrain"
)

type SubscribeCMD struct {
	bot *fatbrain.FatBotBrain
}

func (c *SubscribeCMD) React(chatID int, text string) bool {
	if !strings.HasPrefix(text, "/subscribe") {
		return false
	}

	if err := c.bot.AddChat(chatID); err != nil {
		logrus.WithError(err).Error("Failed adding a new chat")
	}
	return true
}

func New(bot *fatbrain.FatBotBrain) *SubscribeCMD {
	return &SubscribeCMD{
		bot: bot,
	}
}
