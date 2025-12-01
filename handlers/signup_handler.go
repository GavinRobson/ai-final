package handlers

import (
	"ai-final/auth"
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"
)

var signupTmpl = template.Must(template.ParseFiles("templates/auth/signup.html"))

type SignupPageData struct {
	Error string
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := signupTmpl.Execute(w, nil); err != nil {
			http.Error(w, "template error", 500)
		}
	case http.MethodPost:
		userCtx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")

		if password != confirmPassword {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<div class="text-red-500 text-sm text-center">
				Passwords do not match!
				</div>`))
			return
		}

		userID, err := auth.Signup(userCtx, username, password)
		if errors.Is(err, auth.ErrUsernameTaken) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<div class="text-red-500 text-sm text-center">
				Username taken!
				</div>`))
			return
		}
		if err != nil {
			http.Error(w, "Signup failed", http.StatusUnauthorized)
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
