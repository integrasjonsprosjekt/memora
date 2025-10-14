package main

import (
	"fmt"
	"log"

	"memora/internal/config"
	"memora/internal/firebase"
	"memora/internal/router"
	"memora/internal/services"

	"github.com/go-playground/validator/v10"
)

// Initialize configuration
func init() {
	config.Init()
}

// Main entry point of the application
func main() {
	// Initialize Firebase client
	client, app, err := firebase.Init()
	if err != nil {
		log.Panic(err)
	}

	auth, err := firebase.NewFirebaseAuth(app)
	if err != nil {
		log.Panic(err)
	}

	// Initialize validator
	validate := validator.New()

	// Initialize repositories and services
	repos := firebase.NewRepositories(client, auth)
	services := services.NewServices(repos, validate)

	// Set up and run the router
	r := router.New()

	// Register routes with the router
	router.Route(r, services)

	// Start the server
	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}
