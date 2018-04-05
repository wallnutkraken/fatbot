package fatdata

import "database/sql"
import (
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

const (
	testDbSQL         = `SELECT message FROM messages LIMIT 1`
	createMessagesSQL = `CREATE TABLE messages (
	message TEXT NOT NULL
);`

	createChatsSQL = `CREATE TABLE chats (
	chat_id int NOT NULL
);`

	insertMessageSQL = `INSERT INTO messages (message) VALUES (?)`
	insertChatSQL    = `INSERT INTO chats (chat_id) VALUES (?)`
	deleteChatSQL    = `DELETE FROM chats WHERE chat_id = ?`
	getChatsSQL      = `SELECT chat_id FROM chats`
	getMessagesSQL   = `SELECT message FROM messages`
)

func CreateDatabase(connectionString string) error {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil
	}

	_, err = db.Exec(createMessagesSQL)
	if err != nil {
		return err
	}
	_, err = db.Exec(createChatsSQL)
	if err != nil {
		return err
	}

	return nil
}

type Data struct {
	db *sql.DB
}

func Connect(connectionString string) (*Data, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Test if the database is created
	if _, err := db.Exec(testDbSQL); err != nil {
		return nil, err
	}

	return &Data{
		db: db,
	}, nil
}

func (d *Data) Close() error {
	logrus.Info("Ending DB connection")
	return d.db.Close()
}

func (d *Data) SaveMessage(msg string) error {
	_, err := d.db.Exec(insertMessageSQL, msg)
	return err
}

func (d *Data) AddChat(chatID int) error {
	_, err := d.db.Exec(insertChatSQL, chatID)
	return err
}

func (d *Data) RemoveChat(chatID int) error {
	_, err := d.db.Exec(deleteChatSQL, chatID)
	return err
}

func (d *Data) GetChats() ([]int, error) {
	rows, err := d.db.Query(getChatsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (d *Data) GetMessages() ([]string, error) {
	rows, err := d.db.Query(getMessagesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]string, 0)
	for rows.Next() {
		var msg string
		if err := rows.Scan(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
