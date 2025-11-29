// Package handlers takes care of all route handling
package handlers

import (
	"errors"
	"ai-final/auth"
	"context"
	"html/template"
	"net/http"
	"time"
)

var loginTmpl = template.Must(template.ParseFiles("templates/auth/login.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := loginTmpl.Execute(w, nil); err != nil {
			http.Error(w, "template error", 500)
		}
	case http.MethodPost:
		userCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		username := r.FormValue("username")
		password := r.FormValue("password")

		userID, err := auth.Login(userCtx, username, password)
		if errors.Is(err, auth.ErrNoUsersFound) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<div class="text-red-500 text-center">User not found!</div>`))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "user_id",
			Value: userID,
			Path:  "/",
		})

		w.Header().Set("HX-Redirect", "/chat")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	}
}
