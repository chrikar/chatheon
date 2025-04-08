package main

import (
	"log"
	"net/http"

	httpAdapter "github.com/chrikar/chatheon/adapters/http"
	"github.com/chrikar/chatheon/adapters/memory"
	"github.com/chrikar/chatheon/application"
)

func main() {
	messageRepo := memory.NewMessageRepository()
	messageService := application.NewMessageService(messageRepo)

	httpHandler := httpAdapter.NewHandler(messageService)

	log.Println("Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
