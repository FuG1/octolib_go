package handlers

import (
	"net/http"
	"octolib/api/services"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		message := services.GetHelloMessage()
		w.Write([]byte(message))
	case http.MethodPost:
		message := services.PostHelloMessage()
		w.Write([]byte(message))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
