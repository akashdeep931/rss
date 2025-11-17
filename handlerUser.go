package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create user: %s", err.Error()))
		return
	}

	respondWithJSON(w, 201, dbUserToUser(user))
}

func handlerFetchUserByAPIKey(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, 200, dbUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user db.User) {
	dbPosts, err := apiCfg.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't get posts by user: %s", err.Error()))
		return
	}

	posts := []Post{}

	for _, dbPost := range dbPosts {
		posts = append(posts, dbPostToPost(dbPost))
	}

	respondWithJSON(w, 200, posts)
}
