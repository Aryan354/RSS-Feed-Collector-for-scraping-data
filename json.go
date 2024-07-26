package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	// error handling
	if err != nil {
		log.Printf("the error happens to be with the payload as %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

// handle error catching with apis

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("There is a 500 level errror")
	}
	// for marshalling purposes

	type errResponse struct {
		Error string `json:"error'`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
