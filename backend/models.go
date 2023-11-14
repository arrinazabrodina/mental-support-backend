package main

import (
	"fmt"
	"time"
)

type TelegramUser struct {
	CategoryId string `json:"categoryId"`
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}

func telegramUserVisibleName(user TelegramUser) string {
	name := user.FirstName
	if len(user.LastName) != 0 {
		name += fmt.Sprintf(" %s", user.LastName)
	}
	if len(user.LastName) != 0 {
		name += fmt.Sprintf(" (@%s)", user.Username)
	}
	return name
}

func NewTelegramUser(Id int64, Username string, FirstName string, LastName string) TelegramUser {
	user := NewTemplateTelegramUser()
	user.Id = Id
	user.Username = Username
	user.FirstName = FirstName
	user.LastName = LastName
	return user
}
func NewTemplateTelegramUser() TelegramUser {
	var user TelegramUser
	user.CategoryId = "user"
	return user
}

type AdminUser struct {
	CategoryId string `json:"categoryId"`
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	//Password []byte `json:"-"`
}

func NewAdminUser(Id int64, Email string, Name string) AdminUser {
	return AdminUser{
		CategoryId: "admin",
		Id:         Id,
		Email:      Email,
		Name:       Name,
	}
}

func getAdminUser() AdminUser {
	return NewAdminUser(1, "arinazabrodina@knu.ua", "Arina Zabrodina")
}

type Chat struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Attachment struct {
	Id        int `json:"id"`
	MessageId int `json:"messageId"`
}

type Message struct {
	//ChatId    int64  `json:"chatId"`
	//AuthorCategory string        `json:"-"` // admin or user
	//AuthorId       int64         `json:"-"`
	Id          int           `json:"id"`
	Chat        *Chat         `json:"chat"`
	Message     string        `json:"message"`
	MessageId   int           `json:"messageId"` // telegram message id
	AdminAuthor *AdminUser    `json:"-"`
	UserAuthor  *TelegramUser `json:"-"`
	Author      Author        `json:"author"`
	Date        time.Time     `json:"date"`
}

type MessageFlattened struct {
	Id             int
	ChatId         int64
	Message        string
	AuthorId       int64
	AuthorCategory string
	MessageId      int
	Date           string
}

type Author struct {
	CategoryId string `json:"categoryId"`
	Id         int64  `json:"id"`
	Name       string `json:"name"`
}
