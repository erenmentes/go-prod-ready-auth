package main

import (
	"log"
	"net/http"

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
	//config.ConnectDatabase()

	// main router
	router := mux.NewRouter()

	// setup /api routers.
	handler.SetupRoutes(router)

	// start listening
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("something went wrong while listening")
	}
}
