package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var messageChan = getWebMessageChan()

type room struct {

	// clients holds all current clients in this room.
	clients map[*client]bool

	// join is a channel for clients wishing to join the room.
	join chan *client

	// leave is a channel for clients wishing to leave the room.
	leave chan *client

	// forward is a channel that holds incoming messages that should be forwarded to the other clients.
	forward chan []byte
}

// newRoom create a new chat room

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	newTgMessages := *getTelegramMessageChan()
	for {
		//print("Listening")
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			print(msg)
			handleMessageFromWeb(msg)
			//for client := range r.clients {
			//	client.receive <- []byte("Here is a string....")
			//}
		case message := <-newTgMessages:
			for client := range r.clients {
				response := NewMessageResponse{Event: "newMessage", Object: message}
				client.receive <- response.asData()
			}
		}
	}
}

func handleMessageFromWeb(data []byte) {
	var rawMessage MessageFromWeb
	err := json.Unmarshal(data, &rawMessage)
	if err != nil {
		print(err)
	}

	var message Message
	message.Message = rawMessage.Text
	message.Chat = &Chat{Id: int64(rawMessage.Id)}
	admin := getAdminUser()
	message.AdminAuthor = &admin
	message.Author = Author{
		CategoryId: admin.CategoryId,
		Id:         admin.Id,
		Name:       admin.Name,
	}

	*messageChan <- message
}

type MessageFromWeb struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type NewMessageResponse struct {
	Event  string  `json:"event"`
	Object Message `json:"data"`
}

func (r *NewMessageResponse) asData() []byte {
	data, _ := json.MarshalIndent(r, "", " ")
	return data
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
