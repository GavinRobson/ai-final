package handlers

import (
	"ai-final/database"
	"fmt"
	"net/http"
	"strings"
	"html"
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
          class="text-sm w-full p-2 hover:bg-gray-800 rounded-lg cursor-pointer transition-all"
          hx-get="/chat/` + convo.ID + `"
					hx-trigger="click"
          hx-target="#chatArea"
          hx-swap="innerHTML"
          hx-push-url="true"
        >
				`)
				sb.WriteString(html.EscapeString(convo.Title))
				sb.WriteString(`</div>`)
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(sb.String()))
		return
	}

	http.Error(w, "invalid method", http.StatusMethodNotAllowed)
}
