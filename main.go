package main

import (
	"duck-cook-user-ms/api"
	"duck-cook-user-ms/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.Start("3001", db.ConnectMongo())
}
