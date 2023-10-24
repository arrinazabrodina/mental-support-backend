package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func initDb() {

}

func initializeDatabase() bool {
	var err error
	// TODO: (Zabrodina) get user, password, addr from os

	config := mysql.Config{
		User:   "root",
		Passwd: "testtest",
		Net:    "tcp",
		Addr:   "127.0.0.1:3307",
		DBName: "mental_health",
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

func addMessage(message Message) {
	if message.Author.CategoryId == "user" {
		addMessageFromTelegram(message)
		return
	}
	// TODO: add message from admin

}

func addMessageFromTelegram(message Message) bool {
	if message.UserAuthor == nil {
		return false
	}

	user := *message.UserAuthor
	updateOrCreateTelegramUser(user)

	if message.Chat == nil {
		log.Printf("Can't add message without chat")
		return false
	}
	chat := *message.Chat
	updateOrCreateChat(chat)

	createMessage(message)

	return true
}

func createMessage(message Message) {
	_, err := db.Exec("INSERT INTO `Message`(messageId, chatId, message, authorId, authorCategory, date) VALUES (?,?,?,?,?,?);",
		message.MessageId,
		message.Chat.Id,
		message.Message,
		message.Author.Id,
		message.Author.CategoryId,
		message.Date,
	)
	if err != nil {
		log.Printf("Error while creating message: %s", err)
	}
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
