package main

import (
	"fmt"
	"net/http"

	"github.com/akashdeep931/rss/internal/auth"
	"github.com/akashdeep931/rss/internal/db"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Auth error: %s", err.Error()))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, "Couldn't find user")
			return
		}

		handler(w, r, user)
	}
}
