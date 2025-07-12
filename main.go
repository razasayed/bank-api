package main

import (
	"bank-api/db"
	"bank-api/router"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	db.Connect()
	router.InitRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
