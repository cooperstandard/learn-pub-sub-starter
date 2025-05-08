package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("failed to establish connection")
	}
	defer connection.Close()

	fmt.Println("connection successful")

	ch, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to open channel")
	}

	body, err := json.Marshal(routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("failed to marshal payload")
	}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), body)
	if err != nil {
		log.Fatalf("failed to publish json")
	}
	fmt.Println("successfully sent json")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("signal received")
}
