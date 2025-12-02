package handlers

import (
	"ai-final/database"
	"fmt"
	"html"
	"net/http"
	"strings"
)

func ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userIDCookie, err := r.Cookie("user_id")
		if err != nil || userIDCookie.Value == "" {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		userID := userIDCookie.Value
		conversations, err := database.GetConversationsByID(r.Context(), userID)
		if err != nil {
			fmt.Errorf("conversations error: %w", err)
			return
		}

		var sb strings.Builder

		if len(conversations) == 0 {
			sb.WriteString(`
				<div>No Conversations Yet!</div>
				`)
		} else {
			for _, convo := range conversations {
				sb.WriteString(`
  			<div
          class="flex flex-row justify-between items-center text-sm w-full p-2 hover:bg-gray-800 rounded-lg cursor-pointer transition-all"
          hx-post="/api/chat/` + convo.ID + `"
					hx-trigger="click"
          hx-target="#chatArea"
          hx-swap="innerHTML"
        >
				`)
				sb.WriteString(html.EscapeString(convo.Title))
				sb.WriteString(`
					<button 
						class="bg-gray-900 hover:bg-red-500 transition-all rounded-lg cursor-pointer p-1"
						onclick="event.stopPropagation();"
						hx-delete="/conversations/` + convo.ID + `"
						hx-trigger="click"
					>
					<svg 
					xmlns="http://www.w3.org/2000/svg" 
					width="16" 
					height="16" 
					viewBox="0 0 24 24" 
					fill="none" 
					stroke="currentColor" 
					stroke-width="2" 
					stroke-linecap="round" 
					stroke-linejoin="round" 
					>
					<path 
					d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"
					/>
					<path 
					d="M3 6h18"
					/>
					<path 
					d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
					/>
					</svg>

					</button>
					`)
				sb.WriteString(`</div>`)
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(sb.String()))
		return
	}

	if r.Method == http.MethodDelete {
		chatID := strings.TrimPrefix(r.URL.Path, "/conversations/")
		userIDCookie, err := r.Cookie("user_id")
		if err != nil || userIDCookie.Value == "" {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		userID := userIDCookie.Value
		err = database.DeleteConversation(r.Context(), userID, chatID)
		if err != nil {
			http.Error(w, "error deleteing conversation", 500)
			return
		}
		w.Header().Set("HX-Trigger", "refreshConversations")
		return
	}

	http.Error(w, "invalid method", http.StatusMethodNotAllowed)
}
