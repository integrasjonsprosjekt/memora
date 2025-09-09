package main

import (
	"fmt"
	"log"

	// firebase "memora/internal/firebase"

	"memora/internal/config"
	"memora/internal/router"
)

func init() {
	config.Init()
}

func main() {
	r := router.New()

	// firebaseClient, err := firebase.Init()
	// if err != nil {
	// 	log.Panic(err)
	// }

	router.Route(r)

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("Starting HTTP server on %s...\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
