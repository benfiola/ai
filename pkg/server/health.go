package server

import (
	"net/http"

	"github.com/go-chi/render"
)

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	text := "ok"
	err := s.Core.Health(r.Context())
	if err != nil {
		status = http.StatusInternalServerError
		text = "not ok"
	}
	render.Status(r, status)
	render.PlainText(w, r, text)
}
