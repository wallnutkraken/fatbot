package fatbrain

import (
	"math/rand"
	"time"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/telegogo"
)

func (f *FatBotBrain) timedTrigger(closeChan chan bool) {
	now := time.Now().UTC()
	var timeToWait int

	if now.Hour() > 20 || now.Hour() < 10 {
		var hrsToMorning int
		if now.Hour() > 20 {
			hrsToMorning = (20 - now.Hour()) + 10
		} else {
			hrsToMorning = 10 - now.Hour()
		}
		timeToWait = rand.Intn(3600 * hrsToMorning)
	} else {
		for timeToWait < 60*5 {
			timeToWait = rand.Intn(3600)
		}
	}

	select {
	case <-closeChan:
		f.continueMessaging = false
		return
	case <-time.After(time.Duration(timeToWait) * time.Second):
		if f.messageCount < MinMessageCountForMessaging {
			return
		}
		chatsToRemove := make([]int, 0)
		for _, chatId := range f.inChats {
			_, err := f.telegram.SendMessage(TeleGogo.SendMessageArgs{
				ChatID: strconv.Itoa(chatId),
				Text:   f.generate(),
			})

			if err != nil {
				logrus.WithError(err).Infof("Sending message to chat [%d] failed, removing from chats", chatId)
				chatsToRemove = append(chatsToRemove, chatId)
			}
		}
		for _, chatId := range chatsToRemove {
			if err := f.removeChat(chatId); err != nil {
				logrus.WithError(err).Errorf("Failed removing chat [%d]", chatId)
			}
		}
	}
}
