package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/wallnutkraken/fatbot/fatai"
	"github.com/wallnutkraken/fatbot/fatbrain"
	"github.com/wallnutkraken/fatbot/fatctrl"
	"github.com/wallnutkraken/fatbot/fatdata"
	"github.com/wallnutkraken/fatbot/fatplugin"
	"github.com/wallnutkraken/fatbot/fatplugin/commands/saycmd"
	"github.com/wallnutkraken/fatbot/fatplugin/commands/subcmd"
	"github.com/wallnutkraken/fatbot/fatplugin/urlcleaner"
)

const (
	connStr            = "fatbot:fatbot@tcp(tgbot_mysql_1:3306)/fatbot"
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

	chats, err := db.GetChats()
	if err != nil {
		logrus.WithError(err).Fatal("Failed loading chats")
	}

	cleaners := []fatplugin.Cleaner{
		urlcleaner.New(),
	}

	lstmSettings := fatai.LSTMSettings{
		SavePath:  os.Getenv("FATBOT_SAVE_PATH"),
		WordCount: os.Getenv("FATBOT_WORD_COUNT"),
	}
	lstm, err := fatai.New(lstmSettings)
	var startTraining bool
	if err != nil {
		startTraining = true
	}

	brainSettings := fatbrain.FatBotSettings{
		TelegramKey:   os.Getenv("FATBOT_TELEGRAM_TOKEN"),
		RefreshPeriod: time.Second * 2,
		Database:      db,
		Chats:         chats,
		Cleaners:      cleaners,
		FatLSTM:       lstm,
		StartTraining: startTraining,
	}

	brain, err := fatbrain.New(brainSettings)
	if err != nil {
		logrus.WithError(err).Fatal("Failed creating bot")
	}
	reactors := []fatplugin.Reactor{
		saycmd.New(brain),
		subcmd.New(brain),
	}
	brain.AddReactors(reactors)
	brain.Start()

	ctrl := fatctrl.New(":1587", brain)
	if err := ctrl.Start(); err != nil {
		logrus.WithError(err).Error("Control API failed. Stopping brain.")
		brain.StopTraining()
		brain.Stop()
	}
}
