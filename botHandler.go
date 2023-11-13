package main

import (
	"context"
	"encoding/json"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
	"time"
)

var newTgMessage = getTelegramMessageChan()

var tgBot *bot.Bot
var tgCtx context.Context

func startBotHandler() {
	go handleNewMessages()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	tgCtx = ctx
	defer cancel()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	bot, err := bot.New(os.Getenv("TELEGRAM_API_TOKEN"), opts...)
	tgBot = bot
	if err != nil {
		panic(err)
	}
	bot.Start(ctx)
}

func handleNewMessages() {
	print("Entered here")
	messagesChan := *getWebMessageChan()
	for {
		select {
		case message := <-messagesChan:
			sendMessage, err := tgBot.SendMessage(tgCtx, &bot.SendMessageParams{
				ChatID: message.Chat.Id,
				Text:   message.Message,
			})
			if err != nil {
				continue
			}
			var _ Message

			message.Chat = &Chat{
				Id:   message.Chat.Id,
				Name: "",
			}
			message.MessageId = sendMessage.ID
			message.Date = time.Unix(int64(sendMessage.Date), 0)

			addMessage(message)
		}
	}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	message = addMessage(message)

	(*newTgMessage) <- message
	marsh, _ := json.MarshalIndent(message, "", " ")
	print(string(marsh))
}

func (message *Message) asData() []byte {
	data, _ := json.MarshalIndent(message, "", " ")
	return data
}
