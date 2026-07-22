package main

import (
	"log"
	"net/http"

	"github.com/erenmentes/go-prod-ready-auth/internal/config"
	"github.com/erenmentes/go-prod-ready-auth/internal/handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// database connection
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// main router
	router := mux.NewRouter()

	// setup /api routers.
	handler.SetupRoutes(router, db)

	// start listening
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("something went wrong while listening")
	}
}
