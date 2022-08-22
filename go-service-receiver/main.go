package main

import (
	"context"
	"encoding/json"
	"fmt"
	endpoints "go-service-receiver/endpoints"
	model "go-service-receiver/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	smux := http.NewServeMux()
	log := log.New(os.Stdout, "server-receiver", log.LstdFlags)

	receivedEndpoint := endpoints.NewReceivedEndpoint()
	smux.Handle("/received", receivedEndpoint)

	server := http.Server{
		Addr:         ":9091",
		Handler:      smux,
		ErrorLog:     log,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  30 * time.Minute,
	}

	go func() {
		log.Println("Server-receiver has started on port 9091")
		err := server.ListenAndServe()

		if err != nil {
			log.Println("Error while starting server: ", err)
			os.Exit(1)
		}

	}()

	go func() {
		connection, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "topic_test", 0)
		connection.SetReadDeadline(time.Now().Add(time.Second * 10))

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "topic_test",
		})

		//bytes := make([]byte, 1e3)

		for {

			msg, err := reader.ReadMessage(context.Background())

			if err != nil {
				break
			}

			fmt.Println(string(msg.Value))
			entity := &model.Entity{}
			err2 := json.Unmarshal(msg.Value, entity)

			if err2 != nil {
				log.Println("Error while encoding received entity", err2)
			}

			model.ReceiveEntity(entity)
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
