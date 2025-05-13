package main

import (
	"log"
	"net/http"
	"octolib/api"
)

func main() {
	router := api.SetupRoutes()
	log.Println("Server is running on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
