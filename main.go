package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello")
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT has not been configured. Please recheck")
	}

	// define the server
	router := chi.NewRouter()

	// CORS config nfoir the apis

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	//API versioning
	v1Router := chi.NewRouter()
	// get request to check health (i.e. ready)
	v1Router.Get("/ready", handler_ready)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
