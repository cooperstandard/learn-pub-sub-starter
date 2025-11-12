package main

import (
	"fmt"
	"log"

	// "os"
	// "os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	_, _, err = pubsub.DeclareAndBind(connection, "peril_topic", "game_logs", routing.GameLogSlug+".*", 0)
	if err != nil {
		log.Println(err)
		log.Fatalf("failed to create log queue")
	}



	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		if input[0] == "pause" {
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Fatalf("failed to publish json")
			}
			fmt.Println("successfully sent pause")
		} else if input[0] == "resume" {

			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Fatalf("failed to publish json")
			}
			fmt.Println("successfully sent resume")
		} else if input[0] == "quit" {
			fmt.Println("exiting...")
			break
		} else {
			fmt.Println(input[0] + " is not a valid command")
		}
	}
}
