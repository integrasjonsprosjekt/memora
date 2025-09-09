package main

import (
	"fmt"
	"log"
	"net/http"

	// firebase "memora/internal/firebase"

	"memora/internal/config"
	"memora/internal/router"
)

func init() {
	config.Init()
}

func main() {
	handler := router.New()

	// firebaseClient, err := firebase.Init()
	// if err != nil {
	// 	log.Panic(err)
	// }

	log.Printf("Starting HTTP server on %s:%s...\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", config.Host, config.Port),
		handler,
	))
}
