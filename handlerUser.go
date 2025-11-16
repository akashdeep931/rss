package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akashdeep931/rss/internal/auth"
	"github.com/akashdeep931/rss/internal/db"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Body not provided")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create user: %s", err.Error()))
		return
	}

	respondWithJSON(w, 201, dbUserToUser(user))
}

func (apiCfg *apiConfig) handlerFetchUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Auth error: %s", err.Error()))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 404, "Couldn't find user")
		return
	}

	respondWithJSON(w, 200, dbUserToUser(user))
}
