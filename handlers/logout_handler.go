package handlers

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "user_id",
		Value: "",
		Path: "/",
		MaxAge: -1,
	})

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}
