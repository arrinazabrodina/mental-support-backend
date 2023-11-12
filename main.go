package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"text/template"
	"time"
)

var r *room

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.filename))
	})
	t.templ.Execute(w, r)
}

func main() {

	if initializeDatabase() == false {
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/messages", getMessages)
	router.Get("/room", handleWebhook)
	router.Get("/chats", getChats)
	router.Get("/", getTestPage)
	router.Get("/test", test)

	r = newRoom()
	go r.run()
	go startBotHandler()

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type Pagination struct {
	Next          int `json:"next"`
	Previous      int `json:"previous"`
	RecordPerPage int `json:"recordPerPage"`
	CurrentPage   int `json:"currentPage"`
	TotalPage     int `json:"totalPage"`
}

type Page struct {
	Objects  interface{} `json:"objects"`
	Metadata Pagination  `json:"metadata"`
}

type Test struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"dbName"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func test(w http.ResponseWriter, r *http.Request) {
	t := Test{}

	t.Host = os.Getenv("DATABASE_HOST")
	t.Port = os.Getenv("DATABASE_PORT")
	t.DbName = os.Getenv("DATABASE_NAME")
	t.User = os.Getenv("DATABASE_USER")
	t.Password = os.Getenv("DATABASE_PASSWORD")
	//str := os.Environ()

	json.NewEncoder(w).Encode(t)
}

func getTestPage(w http.ResponseWriter, r *http.Request) {
	t := &templateHandler{filename: "chat.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.filename))
	})
	t.templ.Execute(w, r)
}

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	r.ServeHTTP(w, req)
}

func getChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r == nil {
		return
	}

	rows, err := db.Query("SELECT id, name FROM Chat")
	if err != nil {
		makeError(http.StatusInternalServerError, fmt.Sprintf("Unable to chats, error: %s", err), &w)
		return
	}

	defer rows.Close()

	chats := []Chat{}

	for rows.Next() {
		var chat Chat

		if err := rows.Scan(&chat.Id, &chat.Name); err != nil {
			continue
		}
		chats = append(chats, chat)
	}

	json.NewEncoder(w).Encode(chats)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r == nil {
		return
	}
	chatIdRaw := r.URL.Query().Get("chatId")

	if len(chatIdRaw) == 0 {
		makeError(http.StatusBadRequest, "`chatId` is required", &w)
		return
	}
	chatId, err := strconv.Atoi(chatIdRaw)

	if err != nil {
		makeError(http.StatusBadRequest, "`chatId` should be big int", &w)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	limit := 10
	var metadata = Pagination{}
	row := db.QueryRow("SELECT count(id) FROM Message WHERE chatId=?", chatId)
	var recordCount int
	errDB := row.Scan(&recordCount)
	print(errDB)
	pagesCount := recordCount / limit

	remainder := (recordCount % limit)
	if remainder != 0 {
		pagesCount = pagesCount + 1
	}
	metadata.TotalPage = pagesCount

	metadata.CurrentPage = page
	metadata.RecordPerPage = limit

	if page <= 0 {
		metadata.Next = page + 1
	} else if page < metadata.TotalPage {
		metadata.Previous = page - 1
		metadata.Next = page + 1
	} else if page == metadata.TotalPage {
		metadata.Previous = page - 1
		metadata.Next = 0
	}

	offset := limit * (page - 1)
	rows, err := db.Query("SELECT id, chatId, message, authorId, authorCategory, messageId, date FROM Message WHERE chatId=? ORDER BY date DESC LIMIT ? OFFSET ?", chatId, limit, offset)

	if err != nil {
		log.Println(err)
		makeError(http.StatusInternalServerError, "Unable to connect to fetch users from database", &w)
		return
	}
	defer rows.Close()

	messagesFlattened := []MessageFlattened{}

	for rows.Next() {
		var message MessageFlattened

		if err := rows.Scan(&message.Id, &message.ChatId, &message.Message, &message.AuthorId, &message.AuthorCategory, &message.MessageId, &message.Date); err != nil {
			continue
		}
		messagesFlattened = append(messagesFlattened, message)
	}
	messages := []Message{}
	chatsCache := make(map[int64]Chat)

	for _, messageFlattened := range messagesFlattened {
		var message Message
		chat := chatsCache[messageFlattened.ChatId]
		if chat == (Chat{}) {

			chat = *getChatById(messageFlattened.ChatId)
			chatsCache[messageFlattened.ChatId] = chat
		}
		message.Chat = &chat
		if messageFlattened.AuthorCategory == "user" {
			user, err := getTelegramUserById(messageFlattened.AuthorId)
			if err != nil {
				makeError(http.StatusInternalServerError, fmt.Sprintf("Unable to get telegram user, error: %s", err), &w)
				return
			}
			if user == nil {
				makeError(http.StatusInternalServerError, "No user with needed id", &w)
				return
			}
			message.UserAuthor = user
			message.Author = makeAuthorFromTelegramUser(*user)
		} else {
			// TODO: Implement
		}

		message.Id = messageFlattened.Id
		message.MessageId = messageFlattened.MessageId
		message.Message = messageFlattened.Message
		date, err := time.Parse("2006-01-02 15:04:05", messageFlattened.Date)
		if err != nil {
			log.Printf("Unable to parse data: %s\n", err)
		}
		message.Date = date
		messages = append(messages, message)
	}

	pageResult := Page{
		Objects:  messages,
		Metadata: metadata,
	}
	jsonErr := json.NewEncoder(w).Encode(pageResult)
	if jsonErr != nil {
		return
	}
}

func getChatById(id int64) *Chat {
	row := db.QueryRow("SELECT id, name FROM Chat WHERE id=?", id)
	if row == nil {
		return nil
	}
	var chat Chat
	err := row.Scan(&chat.Id, &chat.Name)
	if err != nil {
		return nil
	}
	return &chat
}

func makeError(errorCode int, message string, writer *http.ResponseWriter) {
	errorMessage := ErrorMessage{Code: errorCode, Message: message}
	response := ErrorResponse{Error: errorMessage}
	if writer == nil {
		return
	}
	w := *writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
