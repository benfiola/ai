package server

import (
	"net/http"

	"github.com/benfiola/ai/pkg/core"
	"github.com/go-chi/render"
)

type AuthenticateRequest core.AuthenticateOpts

func (ar *AuthenticateRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	data := AuthenticateRequest{}
	err := render.Bind(r, &data)
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	token, err := s.Core.Authenticate(r.Context(), core.AuthenticateOpts(data))
	if err != nil {
		render.Render(w, r, ErrorBadRequest(err))
		return
	}

	render.PlainText(w, r, token)
}
