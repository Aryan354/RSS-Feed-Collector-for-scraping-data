package main

import "net/http"

func handler_ready(w http.ResponseWriter, r *http.Request) {

	respondWithJson(w, 200, struct{}{})
}
