package main

import (
	"go-service-producer/keycloak"
	"log"
	"os"
	"os/signal"
)

func main() {

	server := keycloak.NewServer("localhost", "9090", keycloak.NewKeycloak())

	go func() {
		log.Println("Server-producer has started on port 9090")
		err := server.Listen()

		if err != nil {
			log.Println("Error while starting server: ", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	server.Shutdown()
}
