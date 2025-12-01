package main

import (
	"ai-final/database"
	"ai-final/handlers"
	"context"
	"log"
	"net/http"
	"time"

	ai "ai-final/openai"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded!")
	}

	_, err := database.InitMongo(ctx)
	if err != nil {
		log.Fatal("DATABASE CONNECTION FAILURE")
	}

	_, err = ai.InitOpenAI(ctx)
	if err != nil {
		log.Fatal("OPENAI CONNECTION FAILURE")
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.HandleFunc("/auth/signup", handlers.SignupHandler)

	http.HandleFunc("/chat", handlers.ChatHandler)

	http.HandleFunc("/conversations", handlers.ConversationsHandler)
	http.HandleFunc("/api/chat/", handlers.PreviousChatHandler)

	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat", http.StatusFound)
	})

	http.HandleFunc("/api/message", handlers.MessageHandler)

	println("Server running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

