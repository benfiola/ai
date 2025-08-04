package server

import (
	"net/http"
	"strings"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		data := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(data) != 2 {
			next.ServeHTTP(w, r)
			return
		}
		token := data[1]
		ctx := s.Core.WithAuth(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
