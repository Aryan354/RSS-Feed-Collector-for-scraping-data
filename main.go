package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	feed, err := urlToFeed("https://feeds.bbci.co.uk/news/world/rss.xml")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

	fmt.Println("Welcome to the RSS server")
	error := godotenv.Load(".env")
	if error != nil {
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

	db_conn := database.New(connStr)
	// defining our connection to dataabase
	apiCfg := apiConfig{
		DB: db_conn,
	}

	// running the concurrent RSS servers
	go startScraping(db_conn, 10, time.Minute)

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

	// new chi router
	v1Router := chi.NewRouter()

	//all the endpoints defined

	v1Router.Get("/ready", handler_ready)
	v1Router.Get("/error", handleError)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUserAPI))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))

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
