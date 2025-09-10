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

	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}
