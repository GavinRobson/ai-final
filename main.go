package main

import (
	"ai-final/database"
	"ai-final/handlers"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type Response struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded!")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY must be set")
	}

	_, err := database.InitMongo(ctx)
	if err != nil {
		log.Fatal("DATABASE CONNECTION FAILURE")
	}

	client := openai.NewClient(apiKey)

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `
			You are a Python tutor. You will answer the user's questions about Python concepts, 
			show very simple implementation examples of python code, point out errors in the user's code,
			create exercises for the user to complete, and give constructive, motivating feedback 
			on the user's code.	This will be a purely textual response so please format 
			accordingly. Please use very short responses that get straight to the point, not 
			a lot of verboseness. Please format your responses in a valid JSON data structure that follows this pattern:
			{
				"title": <If this is the first message you are responding with or you believe the
				title of the converstaion should be updated, create a new or updated title based
				on the context of the converstaion. Else, keep it an empty string "">
				"message": <The text response you will tell the user>,
				"code": <If there is code to be shown to the user, place it here. Else, keep it an empty string "">
			},
			Do not use markdown format. Your full response should be a string that starts with { to show the start of the json,
			and ends with } to show the end of the json. If the code has new lines, please place the escapes for the newline 
			in their proper locations. Also please do the same thing with tabs with the escape for tabs.
			`,
		},
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/auth/login", handlers.LoginHandler)
	http.HandleFunc("/auth/signup", handlers.SignupHandler)

	http.HandleFunc("/chat", handlers.ChatHandler)

	http.HandleFunc("/conversations", handlers.ConversationsHandler)
	http.HandleFunc("/api/chat/", handlers.PreviousChatHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat", http.StatusFound)
	})

	http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		updatedMessages := handleMessage(w, r, client, messages)
		messages = updatedMessages
	})

	println("Server running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

func handleMessage(w http.ResponseWriter, r *http.Request, client *openai.Client, messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	r.ParseForm()
	message := r.FormValue("message")

	respText, updatedMessages := getOpenAIResponse(message, messages, client)

	var parsedResp Response
	if err := json.Unmarshal([]byte(respText), &parsedResp); err != nil {
		log.Fatal(err)
	}

	textMessage := parsedResp.Message
	codeMessage := parsedResp.Code
	codeMessage = strings.ReplaceAll(codeMessage, "\\n", "\n")
	codeMessage = strings.ReplaceAll(codeMessage, "\\t", "\t")
	title := parsedResp.Title

	if codeMessage == "" {
		botMessage := `
		<div class="my-2 bg-gray-700 p-3 rounded-lg self-start max-w-[80%]">
		` + template.HTMLEscapeString(textMessage) + `
		</div>
		`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(botMessage))
		return updatedMessages
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
		database.AddNewConversation(title, updatedMessages, r, w)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(botMessage))
	return updatedMessages
}

func getOpenAIResponse(
	input string,
	messages []openai.ChatCompletionMessage,
	client *openai.Client) (string, []openai.ChatCompletionMessage) {
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: messages,
	})

	if err != nil {
		return "Error: " + err.Error(), nil
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})

	return resp.Choices[0].Message.Content, messages
}
