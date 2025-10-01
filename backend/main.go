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

func init() {
	config.Init()
}

func main() {
	client, err := firebase.Init()
	if err != nil {
		log.Panic(err)
	}

	validate := validator.New()

	repos := firebase.NewRepositories(client)
	services := services.NewServices(repos, validate)

	r := router.New()

	router.Route(r, services)

	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}
