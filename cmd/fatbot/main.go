package main

import (
	"os"
	"strconv"

	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/fatbot/fatbrain"
	"github.com/wallnutkraken/fatbot/fatdata"
	"github.com/wallnutkraken/fatbot/fatplugin"
	"github.com/wallnutkraken/fatbot/fatplugin/commands/saycmd"
	"github.com/wallnutkraken/fatbot/fatplugin/commands/subcmd"
	"github.com/wallnutkraken/fatbot/fatplugin/urlcleaner"
)

const (
	connStr = "fatbot:fatbot@tcp(tgbot_mysql_1:3306)/fatbot"
	defaultChainLength = 1
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	db, err := fatdata.Connect(connStr)
	if err != nil {
		logrus.WithError(err).Info("Failed connecting to DB, attempting to create...")
		if err := fatdata.CreateDatabase(connStr); err != nil {
			logrus.WithError(err).Fatal("Failed creating database")
		}
		db, err = fatdata.Connect(connStr)
		if err != nil {
			logrus.WithError(err).Fatal("Failed connecting to database after creation")
		}
	}
	defer db.Close()

	messages, err := db.GetMessages()
	if err != nil {
		logrus.WithError(err).Fatal("Failed loading messages")
	}
	chats, err := db.GetChats()
	if err != nil {
		logrus.WithError(err).Fatal("Failed loading chats")
	}

	cleaners := []fatplugin.Cleaner{
		urlcleaner.New(),
	}

	brain, err := fatbrain.New(getChainLength(), 8, os.Getenv("FATBOT_TELEGRAM_TOKEN"), 2, db,
		chats, cleaners)
	if err != nil {
		logrus.WithError(err).Fatal("Failed creating bot")
	}
	reactors := []fatplugin.Reactor{
		saycmd.New(brain),
		subcmd.New(brain),
	}
	brain.AddReactors(reactors)

	for _, msg := range messages {
		brain.FeedString(msg)
	}
	brain.Start()

	//TODO
	time.Sleep(time.Hour * 90000)
	brain.Stop()
}

func getChainLength() int {
	chLen := os.Getenv("FATBOT_CHAIN_LENGTH")
	if chLen == "" {
		return defaultChainLength
	}
	chLenInt, err := strconv.Atoi(chLen)
	if err != nil {
		logrus.Errorf("Invalid chain length value [%s]", chLen)
		return defaultChainLength
	}

	return chLenInt
}