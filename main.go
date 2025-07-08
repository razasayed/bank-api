package main

import (
	"bank-api/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	db.Connect()
}
