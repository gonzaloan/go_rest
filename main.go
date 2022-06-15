package main

import (
	"context"
	"goapi/handlers"
	"goapi/middleware"
	"goapi/server"
	"goapi/websocket"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	databaseUrl := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        port,
		JWTSecret:   jwtSecret,
		DatabaseUrl: databaseUrl,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	hub := websocket.NewHub()
	r.Use(middleware.CheckAuthMiddleware(s))
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	r.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/posts", handlers.ListPostHandler(s)).Methods(http.MethodGet)
	go hub.Run()
	r.HandleFunc("/ws", hub.HandleWebSocket)

}
