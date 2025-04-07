package main

import (
	"log"
	"net/http"

	httpAdapter "github.com/chrikar/chatheon/adapters/http"
	"github.com/chrikar/chatheon/adapters/memory"
	"github.com/chrikar/chatheon/adapters/notification"
	"github.com/chrikar/chatheon/application"
)

func main() {
	repo := memory.NewMessageRepository()
	notifier := notification.NewConsoleNotifier()
	service := application.NewChatService(repo, notifier)

	httpHandler := httpAdapter.NewHandler(service)

	log.Println("Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
