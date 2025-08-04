package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/benfiola/ai/pkg/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CreateUserRequest core.CreateUserOpts

func (cur *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	data := CreateUserRequest{}
	err := render.Bind(r, &data)
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	id, err := s.Core.CreateUser(context.Background(), core.CreateUserOpts(data))
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	w.Header().Add("Location", fmt.Sprintf("/user/%d", id))
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	user, err := s.Core.GetUser(r.Context(), id)
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	auth := s.Core.GetAuth(r.Context())
	if auth.User == 0 {
		render.Render(w, r, ErrorBadRequest(fmt.Errorf("invalid auth context")))
		return
	}

	user, err := s.Core.GetUser(r.Context(), auth.User)
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
