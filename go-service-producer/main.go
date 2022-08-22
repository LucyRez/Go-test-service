package main

import (
	"context"
	endpoints "go-service-producer/endpoints"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	smux := http.NewServeMux()
	log := log.New(os.Stdout, "server-producer", log.LstdFlags)

	entityEndpoint := endpoints.NewEntityEndpoint()
	smux.Handle("/entity/", entityEndpoint)

	submitEndpoint := endpoints.NewSubmitEndpoint()
	smux.Handle("/submit", submitEndpoint)

	server := http.Server{
		Addr:         ":9090",
		Handler:      smux,
		ErrorLog:     log,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  30 * time.Minute,
	}

	go func() {
		log.Println("Server-producer has started on port 9090")
		err := server.ListenAndServe()

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

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
