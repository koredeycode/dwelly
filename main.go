package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/koredeycode/dwelly/dwelly-api/handlers"
	"github.com/koredeycode/dwelly/dwelly-api/routes"
	cloudinaryutil "github.com/koredeycode/dwelly/internal/cloudinary"
	"github.com/koredeycode/dwelly/internal/database"
	"github.com/koredeycode/dwelly/internal/redis"
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
	redisClient := redis.InitRedis()
	cloudinaryClient := cloudinaryutil.NewClient()

	apiCfg := handlers.APIConfig{
		DB:         db,
		Redis:      redisClient,
		Cloudinary: cloudinaryClient,
	}

	router := routes.SetUpRouter(&apiCfg)

	routes.PrintRoutes(router)

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
