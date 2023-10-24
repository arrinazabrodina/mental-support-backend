package main

import (
	"context"
	"encoding/json"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func startBotHandler() bool {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	//print(os.Getenv("TELEGRAM_API_TOKEN"))
	bot, err := bot.New(os.Getenv("TELEGRAM_API_TOKEN"), opts...)
	if err != nil {
		panic(err)
		return false
	}
	bot.Start(ctx)
	return true
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   strconv.FormatInt(update.Message.Chat.ID, 10),
	})
	if update.Message.From.IsBot {
		return
	}

	tgUser := NewTelegramUser(
		update.Message.From.ID,
		update.Message.From.Username,
		update.Message.From.FirstName,
		update.Message.From.LastName)

	chatName := update.Message.Chat.Title
	if len(chatName) == 0 {
		chatName = telegramUserVisibleName(tgUser)
	}

	tgChat := Chat{
		Id:   update.Message.Chat.ID,
		Name: chatName,
	}

	author := Author{
		CategoryId: "user",
		Id:         tgUser.Id,
		Name:       telegramUserVisibleName(tgUser),
	}

	message := Message{
		//Id:             0,
		//ChatId:     tgChat.Id,
		Chat:       &tgChat,
		Message:    update.Message.Text,
		MessageId:  update.Message.ID,
		UserAuthor: &tgUser,
		Author:     author,
		Date:       time.Unix(int64(update.Message.Date), 0),
	}

	addMessage(message)
	marsh, _ := json.MarshalIndent(message, "", " ")
	print(string(marsh))
}
