package server

import (
	"log"
	"net/http"

	"github.com/letitloose/user-app/cmd/config"
	"github.com/letitloose/user-app/pkg/user"
)

type Server struct {
	config      *config.Config
	userService *user.UserService
}

func NewServer(config *config.Config, userService *user.UserService) *Server {
	return &Server{
		config:      config,
		userService: userService,
	}
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	log.Println("adding user handlers")
	server.userService.AddHandlersToMux(mux)
	return mux
}

func (server *Server) Run() error {

	log.Println("starting server on port 8080")
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server.Handler(),
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
