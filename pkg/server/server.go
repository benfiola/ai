package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/benfiola/ai/pkg/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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
	server := Server{
		Address: address,
		Core:    core,
		Logger:  logger,
		Router:  router,
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.AllowContentType(("application/json")))
	router.Use(slogchi.New(logger))
	router.Use(server.AuthMiddleware)

	router.Route("/api", func(r chi.Router) {
		r.Post("/authenticate", server.Authenticate)
		r.Get("/health", server.Health)
		r.Post("/user", server.CreateUser)
		r.Get("/user/{id}", server.GetUser)
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

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func (er *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.Status)
	render.JSON(w, r, er)
	return nil
}

func ErrorBadRequest(err error) *ErrorResponse {
	return &ErrorResponse{
		Error:  err.Error(),
		Status: 400,
	}
}

func ErrorInternalServer(err error) *ErrorResponse {
	return &ErrorResponse{
		Error:  err.Error(),
		Status: 500,
	}
}
