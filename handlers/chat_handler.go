package handlers

import (
	"ai-final/database"
	"context"
	"html/template"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var chatShellTmpl = template.Must(template.ParseFiles("templates/chat/chat.html"))

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/chat")

	userIDCookie, err := r.Cookie("user_id")
	if err != nil || userIDCookie.Value == "" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	userID := userIDCookie.Value

	userCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	client, err := database.InitMongo(userCtx)
	if err != nil {
		http.Error(w, "database error", 500)
		return
	}

	content := `<div class="text-gray-300">
		<h1 class="text-2xl font-bold mb-2">Welcome!</h1>
		<p>Select a conversation on the left or start a new one.</p>
	</div>`

	if path != "" && path != "/" {
		chatID := strings.TrimPrefix(path, "/")

		oid, err := bson.ObjectIDFromHex(chatID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var doc bson.M
		err = client.Collection("conversations").FindOne(r.Context(), bson.M{
			"_id":    oid,
			"userId": userID,
		}).
			Decode(&doc)

		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return
		}

		if err != nil {
			http.Error(w, "internal db error", 500)
			return
		}

		chatTitle, _ := doc["title"].(string)

		content = `<div class="p-4">
			<button 
				hx-get="/chat" 
				hx-target="#chatArea" 
				hx-push-url="true"
				class="mb-4 px-3 py-1 border border-gray-700 rounded"
			>
				← Back
			</button>
			<h1 class="text-xl font-semibold mb-2">` + chatTitle + `</h1>
			<p class="text-gray-400 text-sm mb-4">Chat ID: ` + oid.Hex() + `</p>
			<div>Chat messages go here...</div>
		</div>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := chatShellTmpl.Execute(w, map[string]any{"Content": template.HTML(content)}); err != nil {
		http.Error(w, "template error", 500)
	}
}
