package main

import (
	"fmt"
	"net/http"

	"github.com/Aryan354.RssServer/internal/auth"
	"github.com/Aryan354.RssServer/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// middlware to get user through API key (authenticated user)
func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api_key, error := auth.GetAPIKey(r.Header)
		if error != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", error))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), api_key)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Couldn't get user %v", err))
			return
		}

		handler(w, r, user)
	}
}
