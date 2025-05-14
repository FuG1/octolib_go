package main

import (
	"log"
	"net/http"
	"octolib/api"
	"octolib/db"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	router := api.SetupRoutes()
	log.Println("Server is running on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
