package main

var newTelegramMessage = make(chan Message)
var newWebMessage = make(chan Message)

func getTelegramMessageChan() *(chan Message) {
	return &newTelegramMessage
}
func getWebMessageChan() *(chan Message) {
	return &newWebMessage
}
