package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/SaiThihan/go-basic/internal/store"
	"github.com/SaiThihan/go-basic/internal/tokens"
	"github.com/SaiThihan/go-basic/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("User not found in context")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.GuestUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Payload{"error": "Invalid auth header"})
			return
		}

		Token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, Token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Payload{"error": "Invalid token"})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Payload{"error": "Invalid token"})
			return
		}

		r = SetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
