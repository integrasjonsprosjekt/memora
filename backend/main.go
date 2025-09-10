package main

import (
	"fmt"
	"log"

	"memora/internal/config"
	"memora/internal/router"
)

func init() {
	config.Init()
}

func main() {
	r := router.New()

	router.Route(r)

	if err := r.Run(fmt.Sprintf("%s:%s", config.Host, config.Port)); err != nil {
		log.Fatal(err)
	}
}
