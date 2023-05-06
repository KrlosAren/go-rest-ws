package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go/rest-ws/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to the Home Server",
			Status:  http.StatusOK,
		})

	}
}
