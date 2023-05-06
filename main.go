package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go/rest-ws/handlers"
	middleware "github.com/go/rest-ws/middlewares"
	"github.com/go/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading environment ")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseURL: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {

	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.CheckAuthMiddleware(s))
	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostDHandler(s)).Methods(http.MethodDelete)

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.ListPosts(s)).Methods(http.MethodGet)

	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
