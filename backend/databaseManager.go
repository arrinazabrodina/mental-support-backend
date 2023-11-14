package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func initDb() {

}

func initializeDatabase() bool {
	var err error
	// TODO: (Zabrodina) get user, password, addr from os
	//DATABASE_HOST
	//DATABASE_PORT
	//DATABASE_NAME
	//DATABASE_USER
	//DATABASE_PASSWORD
	//TELEGRAM_API_TOKEN

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	dbName := os.Getenv("DATABASE_NAME")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	//panic("Connecting addr: " + addr)
	//return true
	config := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbName,
	}
	db, err = sql.Open("mysql", config.FormatDSN())

	if err != nil {
		fmt.Printf("Unable to open connection with database\nError: %s\n", err)
		return false
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Database connected\n")
	return true
}

func getTelegramUserById(id int64) (*TelegramUser, error) {
	row := db.QueryRow("SELECT id, username, firstName, lastName FROM TelegramUser WHERE id=?", id)
	if row == nil {
		return nil, nil
	}
	user := NewTemplateTelegramUser()
	err := row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func makeAuthorFromTelegramUser(user TelegramUser) Author {
	return Author{
		CategoryId: user.CategoryId,
		Id:         user.Id,
		Name:       telegramUserVisibleName(user),
	}
}

func addMessage(message Message) Message {
	if message.Author.CategoryId == "user" {
		return addMessageFromTelegram(message)
	}
	return addMessageFromWeb(message)
	//return message
}

func addMessageFromWeb(message Message) Message {
	id := createMessage(message)
	message.Id = int(id)
	return message
}

func addMessageFromTelegram(message Message) Message {
	if message.UserAuthor == nil {
		return message
	}

	user := *message.UserAuthor
	updateOrCreateTelegramUser(user)

	if message.Chat == nil {
		log.Printf("Can't add message without chat")
		return message
	}
	chat := *message.Chat
	updateOrCreateChat(chat)

	id := createMessage(message)
	message.Id = int(id)

	return message
}

func createMessage(message Message) int64 {
	result, err := db.Exec("INSERT INTO `Message`(messageId, chatId, message, authorId, authorCategory, date) VALUES (?,?,?,?,?,?);",
		message.MessageId,
		message.Chat.Id,
		message.Message,
		message.Author.Id,
		message.Author.CategoryId,
		message.Date,
	)

	if err != nil {
		return 0
		log.Printf("Error while creating message: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0
		//return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	return id
}

func updateOrCreateChat(chat Chat) {
	dbUser := fetchChat(chat.Id)
	if dbUser == nil {
		createChat(chat)
		return
	}

	_, err := db.Exec("UPDATE `Chat` SET name=? WHERE id=?",
		chat.Name,
		chat.Id)

	if err != nil {
		log.Printf("Failed to update chat: %s\n", err)
	}
}

func fetchChat(id int64) *Chat {
	if db == nil {
		return nil
	}
	row := db.QueryRow("SELECT id, name FROM `Chat` WHERE id=?", id)

	if row == nil {
		return nil
	}
	var chat Chat
	err := row.Scan(&chat.Id, &chat.Name)
	if err != nil {
		log.Printf("Failed to scan user: %s\n", err)
		return nil
	}

	return &chat
}

func createChat(chat Chat) {
	_, err := db.Exec("INSERT INTO `Chat`(id, name) VALUES (?,?);",
		chat.Id,
		chat.Name)
	if err != nil {
		log.Printf("Failed to create chat: %s\n", err)
	}
}

func updateOrCreateTelegramUser(user TelegramUser) {
	dbUser := fetchTelegramUser(user.Id)
	if dbUser == nil {
		createTelegramUser(user)
		return
	}

	_, err := db.Exec("UPDATE `TelegramUser` SET username=?,firstName=?,lastName=? WHERE id=?",
		user.Username,
		user.FirstName,
		user.LastName,
		user.Id)

	if err != nil {
		log.Printf("Failed to update user: %s\n", err)
	}
}

func fetchTelegramUser(id int64) *TelegramUser {
	if db == nil {
		return nil
	}
	row := db.QueryRow("SELECT id, username, firstName, lastName FROM `TelegramUser` WHERE id=?", id)

	if row == nil {
		return nil
	}
	user := NewTemplateTelegramUser()
	err := row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName)
	if err != nil {
		log.Printf("Failed to scan user: %s\n", err)
		return nil
	}

	return &user
}

func createTelegramUser(user TelegramUser) {
	_, err := db.Exec("INSERT INTO `TelegramUser`(id, username, firstName, lastName) VALUES (?,?,?,?);",
		user.Id,
		user.Username,
		user.FirstName,
		user.LastName)
	if err != nil {
		log.Printf("Failed to create user: %s\n", err)
	}
}
