package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	handler "github.com/chrikar/chatheon/adapters/http"
	"github.com/chrikar/chatheon/adapters/memory"
	"github.com/chrikar/chatheon/application"
	"github.com/chrikar/chatheon/internal/auth"
)

func main() {
	fmt.Printf("Chatheon server started at %s\n", time.Now().Format(time.RFC1123))

	messageRepo := memory.NewMessageRepository()
	messageService := application.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageService)

	jwtManager := auth.NewJWTManager("your-secret-key", time.Hour)

	// Repositories
	userRepo := memory.NewUserRepository()
	convRepo := memory.NewConversationRepository()

	// Services
	userService := application.NewUserService(userRepo, jwtManager)
	convService := application.NewConversationService(convRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	convHandler := handler.NewConversationHandler(convService)

	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", userHandler.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.LoginUser).Methods(http.MethodPost)

	// Protected routes
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(auth.JWTMiddleware(jwtManager))

	// Conversation endpoints
	secured.HandleFunc("/conversations", convHandler.CreateConversation).Methods(http.MethodPost)
	secured.HandleFunc("/conversations", convHandler.GetConversations).Methods(http.MethodGet)

	secured.HandleFunc("/messages", messageHandler.CreateMessage).Methods(http.MethodPost)
	secured.HandleFunc("/messages", messageHandler.GetMessages).Methods(http.MethodGet)
	secured.HandleFunc("/messages/{id}/status", messageHandler.UpdateStatus).Methods(http.MethodPut)

	log.Println("Chat server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
