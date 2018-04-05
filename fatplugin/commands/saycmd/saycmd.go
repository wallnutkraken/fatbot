package saycmd

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/fatbot/fatbrain"
)

type SayCMD struct {
	bot *fatbrain.FatBotBrain
}

func (c *SayCMD) React(chatID int, text string) bool {
	if !strings.HasPrefix(text, "/say") {
		return false
	}

	if err := c.bot.SendMessage(chatID); err != nil {
		logrus.WithError(err).Errorf("Failed responding to /say command in chat [%d] with text [%s]",
			chatID, text)
	}
	return true
}

func New(bot *fatbrain.FatBotBrain) *SayCMD {
	return &SayCMD{
		bot: bot,
	}
}
