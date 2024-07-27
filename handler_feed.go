package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Aryan354.RssServer/internal/database"
	"github.com/google/uuid"
)

//fuction to create new user

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Print("There has been a client-side error")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeed(r.Context())

	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "Returned no feed")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
