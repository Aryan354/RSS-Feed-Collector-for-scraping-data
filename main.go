package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Aryan354.RssServer/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Welcome to the RSS server")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	fmt.Printf("The RSS server has initiated on port: %v", port)
	if port == "" {
		log.Fatal("PORT has not been configured. Please recheck")
	}

	// Database connection
	db := os.Getenv("DB_URL")
	if db == "" {
		log.Fatal("Database connection string is missing")
	}

	connStr, err := sql.Open("postgres", db)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer connStr.Close()

	// defining our connection to dataabase
	apiCfg := apiConfig{
		DB: database.New(connStr),
	}

	// Define the server
	router := chi.NewRouter()

	// CORS configuration for the APIs
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	// API versioning
	v1Router := chi.NewRouter()
	//all the routers defined

	v1Router.Get("/ready", handler_ready)
	v1Router.Get("/error", handleError)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUserAPI))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)

	//mounting api versoning
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
