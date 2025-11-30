package handlers

import (
	"ai-final/database"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var chatShellTmpl = template.Must(template.ParseFiles("templates/chat/chat.html"))

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	userIDCookie, err := r.Cookie("user_id")
	if err != nil || userIDCookie.Value == "" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := chatShellTmpl.Execute(w, nil); err != nil {
		http.Error(w, "template error", 500)
	}
}

type Response struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type StoredMessage struct {
	Role string `bson:"role" json:"role"`
	Content string `bson:"content" json:"content"`
}

type ConversationDoc struct {
	ID bson.ObjectID `bson:"_id"`
	UserID string `bson:"userId"`
	Messages []StoredMessage `bson:"messages"`
}

func PreviousChatHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/chat")

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

	chatID := strings.TrimPrefix(path, "/")

	oid, err := bson.ObjectIDFromHex(chatID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var convo ConversationDoc
	err = client.Collection("conversations").FindOne(r.Context(), bson.M{
		"_id":    oid,
		"userId": userID,
	}).Decode(&convo)

	if err == mongo.ErrNoDocuments {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, "internal db error", 500)
		return
	}

	messages := convo.Messages 

	var sb strings.Builder

	for _, message := range messages {

		if message.Role == openai.ChatMessageRoleUser {
			sb.WriteString(`<div class="my-2 bg-blue-600 p-3 rounded-lg self-end max-w-[80%]">
			` + message.Content + `</div>`)
		}

		if message.Role == openai.ChatMessageRoleAssistant {
			var parsedResp Response
			if err := json.Unmarshal([]byte(message.Content), &parsedResp); err != nil {
				http.Error(w, "parsing error", 500)
				return
			}

			textMessage := parsedResp.Message
			codeMessage := parsedResp.Code
			codeMessage = strings.ReplaceAll(codeMessage, "\\n", "\n")
			codeMessage = strings.ReplaceAll(codeMessage, "\\t", "\t")

			if (codeMessage == "") {
				sb.WriteString(`
					<div class="my-2 bg-gray-700 p-3 rounded-lg self-start max-w-[80%]">
					` + template.HTMLEscapeString(textMessage) + `</div>`)
			} else {
				sb.WriteString(`
					<div class="my-2 bg-gray-700 p-2 rounded-lg self-start flex flex-col max-w-[80%]">
					<div class="p-3 rounded-md">
					` + template.HTMLEscapeString(textMessage) + `
					</div>

					<div class="p-3 font-bold rounded-md">
					<pre class="bg-gray-900 text-gray-100 p-3 rounded-lg whitespace-pre-wrap tab-size-2">` + template.HTMLEscapeString(codeMessage) + `</pre>
					</div>
					</div>
					`)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(sb.String()))
}
