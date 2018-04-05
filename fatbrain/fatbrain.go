package fatbrain

import (
	"errors"
	"fmt"
	"strings"

	"io"

	"bytes"

	"time"

	"strconv"

	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/exercism/arkov/chain"
	"github.com/wallnutkraken/fatbot/fatdata"
	"github.com/wallnutkraken/fatbot/fatplugin"
	"github.com/wallnutkraken/telegogo"
)

const (
	MinChainLength              = 1
	MaxChainLength              = 3
	MinMessageCountForMessaging = 100
)

var ErrInvalidLength = errors.New(fmt.Sprintf("The chain length MUST be between %d and %d.",
	MinChainLength, MaxChainLength))

type FatBotBrain struct {
	chain             *chain.Chain
	maxWordLength     int
	telegram          TeleGogo.Client
	refreshSeconds    int
	lastID            int
	inChats           []int
	database          *fatdata.Data
	messageCount      int
	messagingChannel  chan bool
	listeningChannel  chan bool
	continueMessaging bool
	reactors          []fatplugin.Reactor
	cleaners          []fatplugin.Cleaner
	chatMutex         *sync.Mutex
}

// New creates a new instance of FatBotBrain
func New(chainLength int, maxWords int, botKey string, refreshSecs int, db *fatdata.Data,
	chats []int, cleaners []fatplugin.Cleaner) (*FatBotBrain, error) {

	if chainLength < MinChainLength || chainLength > MaxChainLength {
		return nil, ErrInvalidLength
	}

	c := chain.NewChain(chainLength)
	bot, err := TeleGogo.NewBot(botKey)
	if err != nil {
		return nil, err
	}
	brain := &FatBotBrain{
		chain:             c,
		maxWordLength:     maxWords,
		telegram:          bot,
		refreshSeconds:    refreshSecs,
		inChats:           chats,
		database:          db,
		messageCount:      0,
		continueMessaging: true,
		reactors:          make([]fatplugin.Reactor, 0),
		cleaners:          cleaners,
		chatMutex:         &sync.Mutex{},
	}

	return brain, nil
}

func (f *FatBotBrain) AddReactors(reactors []fatplugin.Reactor) {
	f.reactors = append(f.reactors, reactors...)
}

func (f *FatBotBrain) Feed(r io.Reader) {
	f.messageCount++
	f.chain.Build(r)
}

func (f *FatBotBrain) FeedString(text string) {
	buf := bytes.NewBufferString(text)
	f.Feed(buf)
}

func (f *FatBotBrain) generate() string {
	text := f.chain.Generate()
	logrus.Infof("Generated message [%s] with [%d] newlines", text, strings.Count(text, "\n"))
	return strings.Split(text, "\n")[0]
}

func (f *FatBotBrain) AddChat(chatID int) error {
	for _, existingChatID := range f.inChats {
		if existingChatID == chatID {
			return errors.New("Chat already added")
		}
	}

	f.chatMutex.Lock()
	f.inChats = append(f.inChats, chatID)
	err := f.database.AddChat(chatID)
	f.chatMutex.Unlock()

	return err
}

func (f *FatBotBrain) removeChat(chatID int) error {
	var chatIndex int
	for index, id := range f.inChats {
		if id == chatID {
			chatIndex = index
		}
	}
	if chatIndex == 0 {
		return errors.New("Could not find that chat")
	}
	f.chatMutex.Lock()
	f.inChats = append(f.inChats[:chatIndex], f.inChats[chatIndex:]...)
	err := f.database.RemoveChat(chatID)
	f.chatMutex.Unlock()

	return err
}

func (f *FatBotBrain) Start() {
	f.listeningChannel = f.startListening()
	f.messagingChannel = f.startMessaging()
}

func (f *FatBotBrain) Stop() {
	f.listeningChannel <- true
	f.messagingChannel <- true
}

func (f *FatBotBrain) startMessaging() chan bool {
	ch := make(chan bool, 0)
	go func(f *FatBotBrain) {
		for f.continueMessaging {
			f.timedTrigger(ch)
		}
	}(f)
	return ch
}

func (f *FatBotBrain) SendMessage(chatID int) error {
	msgText := f.generate()
	logrus.Infof("Sending message to [%d]: [%s]", chatID, msgText)
	_, err := f.telegram.SendMessage(TeleGogo.SendMessageArgs{
		ChatID: strconv.Itoa(chatID),
		Text:   msgText,
	})

	return err
}

func (f *FatBotBrain) startListening() chan bool {
	ch := make(chan bool, 0)
	go func(ch chan bool, brain *FatBotBrain) {
		for {
			select {
			case <-ch:
				return
			case <-time.After(time.Second * time.Duration(f.refreshSeconds)):
				updates, err := f.telegram.GetUpdates(TeleGogo.GetUpdatesOptions{Offset: f.lastID + 1})
				if err != nil {
					logrus.WithError(err).Error("Failed getting updates")
					continue
				}

				msgsToSave := make([]TeleGogo.Update, 0)
				for _, update := range updates {
					if update.Message.Text != "" {
						logrus.Infof("Got message [%s]", update.Message.Text)
						cleanText := update.Message.Text
						for _, cleaner := range f.cleaners {
							cleanText = cleaner.Clean(cleanText)
						}

						f.FeedString(cleanText)

						var reacted bool
						for _, reactor := range f.reactors {
							if reacted = reactor.React(update.Message.Chat.ID, cleanText); reacted {
								// The bot reacted, continue
								break
							}
						}
						if !reacted {
							msgsToSave = append(msgsToSave, update)
						}
					} else {
						continue
					}
				}
				if len(updates) > 0 {
					f.lastID = updates[len(updates)-1].ID
					go f.saveMesages(msgsToSave)
				}
			}
		}
	}(ch, f)

	return ch
}

func (f *FatBotBrain) saveMesages(updates []TeleGogo.Update) {
	for _, update := range updates {
		cleanText := update.Message.Text
		for _, cleaner := range f.cleaners {
			cleanText = cleaner.Clean(cleanText)
		}
		if err := f.database.SaveMessage(cleanText); err != nil {
			logrus.WithError(err).Errorf("Failed saving message [%s] to database", cleanText)
		}
	}
}
