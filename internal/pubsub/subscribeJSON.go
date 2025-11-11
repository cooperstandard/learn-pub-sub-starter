package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" 0 or "transient" 1
	handler func(T) AckType, // acktype
) error {
	c, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	receiveCh, err := c.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(receiveCh <-chan amqp.Delivery) {
		for delivery := range receiveCh {
			var message T
			// fmt.Printf("Raw body: %s\n", string(delivery.Body))
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			ackType := handler(message)
			switch ackType {
			case Ack:
				delivery.Ack(false)
			case NackDiscard:
				delivery.Nack(false, false)
			case NackRequeue:
				delivery.Nack(false, true)
			default:
				delivery.Ack(false)
			}
		}
	}(receiveCh)

	return nil
}
