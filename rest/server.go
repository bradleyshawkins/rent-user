package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/bradleyshawkins/rent-user/identity"

	"github.com/go-chi/chi"
)

type Server struct {
	Mux           *chi.Mux
	signUpManager *identity.SignUpManager
}

func NewServer(signUpManager *identity.SignUpManager) *Server {
	log.Println("Creating Router")

	mux := chi.NewRouter()

	s := &Server{
		Mux:           mux,
		signUpManager: signUpManager,
	}

	mux.Post("/users", ErrorHandler(s.RegisterUser))

	return s
}

func (s *Server) Start(port string) func(ctx context.Context) error {
	srv := http.Server{
		Addr:    ":" + port,
		Handler: s.Mux,
	}

	go func() {
		log.Println("Starting http server ...")

		err := http.ListenAndServe(":"+port, s.Mux)
		if err != nil {
			log.Println("Error shutting down server. Error:", err)
		}
	}()

	return func(ctx context.Context) error {
		log.Println("Shutting down http server ...")
		return srv.Shutdown(ctx)
	}
}
