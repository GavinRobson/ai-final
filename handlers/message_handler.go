package handlers

import (
	"ai-final/database"
	ai "ai-final/openai"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type StoredMessages struct {
	ChatID   string
	Messages []openai.ChatCompletionMessage
}

var persistentMessages = make(map[string]StoredMessages)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	userIDCookie, err := r.Cookie("user_id")
	if err != nil || userIDCookie.Value == "" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	userID := userIDCookie.Value

	r.ParseForm()

	input := r.FormValue("message")

	respText, updatedMessages, err := ai.GetOpenAIResponse(input, persistentMessages[userID].Messages)
	if err != nil {
		http.Error(w, "error getting openai response", 500)
	}

	var parsedResp Response
	if err := json.Unmarshal([]byte(respText), &parsedResp); err != nil {
		log.Fatal(err)
	}

	textMessage := parsedResp.Message
	codeMessage := parsedResp.Code
	codeMessage = strings.ReplaceAll(codeMessage, "\\n", "\n")
	codeMessage = strings.ReplaceAll(codeMessage, "\\t", "\t")
	title := parsedResp.Title

	if persistentMessages[userID].Messages == nil {
		current := persistentMessages[userID]
		current.Messages = updatedMessages
		persistentMessages[userID] = current
		chatID, err := database.AddNewConversation(title, userID, persistentMessages[userID].Messages)
		if err != nil {
			http.Error(w, "error adding new conversation", 500)
			return
		}
		current.ChatID = chatID
		persistentMessages[userID] = current
	} else {
		current := persistentMessages[userID]
		current.Messages = updatedMessages
		persistentMessages[userID] = current
		userMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		}
		botMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: respText,
		}
		database.AddMessageToConversation(r.Context(), title, persistentMessages[userID].ChatID, userID, userMessage, botMessage)
	}

	if codeMessage == "" {
		botMessage := `
		<div class="my-2 bg-gray-700 p-3 rounded-lg self-start max-w-[80%]">
		` + template.HTMLEscapeString(textMessage) + `
		</div>
		`
		w.Header().Set("HX-Trigger", "refreshConversations")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(botMessage))
		return
	}
	botMessage := `
	<div class="my-2 bg-gray-700 p-2 rounded-lg self-start flex flex-col max-w-[80%]">
	<div class="p-3 rounded-md">
	` + template.HTMLEscapeString(textMessage) + `
	</div>

	<div class="p-3 font-bold rounded-md">
	<pre class="bg-gray-900 text-gray-100 p-3 rounded-lg whitespace-pre-wrap tab-size-2">` + template.HTMLEscapeString(codeMessage) + `</pre>
	</div>
	</div>
	`

	w.Header().Set("HX-Trigger", "refreshConversations")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(botMessage))
}

func GetMessages(userID string) []openai.ChatCompletionMessage {
	return persistentMessages[userID].Messages
}

func SetMessages(chatID, userID string, messages []openai.ChatCompletionMessage) {
	if _, exists := persistentMessages[userID]; exists {
		current := persistentMessages[userID]
		current.ChatID = chatID
		current.Messages = messages
		persistentMessages[userID] = current
		return
	}
	persistentMessages[userID] = StoredMessages{
		ChatID:   chatID,
		Messages: messages,
	}
}

func AddMessage(userID string, message openai.ChatCompletionMessage) {
	current := persistentMessages[userID].Messages
	updated := append(current, message)
	current1 := persistentMessages[userID]
	current1.Messages = updated
	persistentMessages[userID] = current1
}
