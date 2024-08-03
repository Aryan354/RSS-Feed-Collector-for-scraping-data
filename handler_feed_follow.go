package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Aryan354.RssServer/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//fuction to create new user

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Print("There has been a client-side error")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "couldn't get feed_follows	")
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "You didn't prove a feedFollwID")
		return
	}

	error := apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if error != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete teh requested feed_follow item from db")
		return
	}
	respondWithJson(w, 200, struct{}{})
}
