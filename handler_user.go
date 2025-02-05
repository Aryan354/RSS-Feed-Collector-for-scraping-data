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

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Print("There has been a client-side error")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		log.Print("There has been a server-side error")
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJson(w, http.StatusOK, UserReturner(user))
}

//function to handle authenticated user

func (apiCfg *apiConfig) handleGetUserAPI(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJson(w, 200, UserReturner(user))
}
