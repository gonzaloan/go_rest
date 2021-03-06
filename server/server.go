package server

import (
	"context"
	"errors"
	"goapi/database"
	"goapi/repository"
	"goapi/websocket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Config has elements that we need to connect
type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

// Server : We need a Config component to create a Server.
type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

// Broker : Component that will handle other components
type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

// Config : Make Broker satisfy Config, with a Receiver Function
func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}
	return &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}, nil
}

// Start : Adds to broker to allow start. And Start API Rest.
func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
