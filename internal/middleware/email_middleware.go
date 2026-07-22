package middleware

import (
	"net/http"
)

func NewEmailMiddleware() *emailMiddleware {
	return &emailMiddleware{}
}

type emailMiddleware struct{}

func (m *emailMiddleware) Apply(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if !ok {
			respondError(w, http.StatusUnauthorized, "unauthenticated")
			return
		}

		if user.IsEmailVerified == nil || !*user.IsEmailVerified {
			respondError(w, http.StatusForbidden, "email not verified")
			return
		}

		next(w, r)
	}
}
