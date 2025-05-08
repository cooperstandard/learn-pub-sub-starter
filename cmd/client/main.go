package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("failed to establish connection")
	}
	defer connection.Close()

	fmt.Println("connection successful")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("unable to welcome client")
	}

	_, _, err = pubsub.DeclareAndBind(connection, "peril_direct", fmt.Sprintf("pause.%s", username), "pause", 1)
	if err != nil {
		log.Fatalf("unable to bind q and channel")
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("signal received")
}
