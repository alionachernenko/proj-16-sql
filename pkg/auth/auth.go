package auth

import (
	"net/http"
	"tasks/internal/database"

	"github.com/rs/zerolog/log"
)

type Auth struct {
	S *database.PostgresStorage
}

func (a *Auth) CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := a.S.GetUser(username)

		if err != nil {
			log.Error().Err(err).Msg("Failed to retrieve user from database")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if user.Password != password {
			log.Error().Err(err).Msg("Unauthorized")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
