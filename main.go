package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/koredeycode/dwelly/api"
	"github.com/koredeycode/dwelly/internal/database"
	"github.com/koredeycode/dwelly/routes"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbURL := os.Getenv("DB_URL")
	if portString == "" {
		log.Fatal("DBURL environment variable not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	db := database.New(conn)
	apiCfg := api.APIConfig{
		DB: db,
	}

	router := routes.SetUpRouter(&apiCfg)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Starting server on port %s\n", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
