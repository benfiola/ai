package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/benfiola/ai/pkg/core"
	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

type Server struct {
	Address string
	Core    *core.Core
	Logger  *slog.Logger
	Router  *chi.Mux
}

type Opts struct {
	Address string
	Core    *core.Core
	Logger  *slog.Logger
}

func New(opts Opts) (*Server, error) {
	address := opts.Address
	if address == "" {
		address = "0.0.0.0:8080"
	}

	core := opts.Core
	if core == nil {
		return nil, fmt.Errorf("core is nil")
	}

	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	router := chi.NewRouter()
	router.Use(slogchi.New(logger))

	server := Server{
		Address: address,
		Core:    core,
		Logger:  logger,
		Router:  router,
	}

	router.Route("/api", func(r chi.Router) {
		r.Get("/health", server.Health)
	})
	return &server, nil
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error)
	server := http.Server{
		Addr:    s.Address,
		Handler: s.Router,
	}

	go func() {
		s.Logger.Info("server starting", "address", s.Address)
		errChan <- server.ListenAndServe()
	}()

	var err error
	select {
	case serr := <-errChan:
		err = serr
	case <-ctx.Done():
		s.Logger.Info("server stopping")
		server.Shutdown(context.Background())
	}

	return err
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	code := http.StatusOK
	text := "ok"
	err := s.Core.Health()
	if err != nil {
		code = http.StatusInternalServerError
		text = "not ok"
	}
	w.WriteHeader(code)
	w.Header().Add("content-type", "text/plain")
	w.Write([]byte(text))
}
