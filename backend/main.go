package main

import (
	"fmt"
	"log"

	"memora/internal/config"
	"memora/internal/firebase"
	"memora/internal/router"
	"memora/internal/services"
)

func init() {
	config.Init()
}

func main() {
	client, err := firebase.Init()
	
	if err != nil {
		log.Panic(err)
	}

	userRepo := firebase.NewFirestoreUserRepo(client)
	userService := services.NewUserService(userRepo)

	r := router.New()

	router.Route(r, userService)

	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}
