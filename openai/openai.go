// Package openai handles ai requests from the application
package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var client *openai.Client

const initialPrompt string =`
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
			`

func InitOpenAI(ctx context.Context) (*openai.Client, error) {
	if client != nil {
		return client, nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	client = openai.NewClient(apiKey)
	return client, nil
}

func GetOpenAIResponse(input string, messages []openai.ChatCompletionMessage) (string, []openai.ChatCompletionMessage, error) {
	if messages == nil {
		messages = append(messages, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			Content: initialPrompt,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: messages,
	})

	if err != nil {
		return "", nil, fmt.Errorf("error getting openai resopnse")
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})

	return resp.Choices[0].Message.Content, messages, nil
}

