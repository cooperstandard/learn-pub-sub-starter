package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int
type SimpleQueueType int

const (
	Ack AckType = iota
	NackDiscard
	NackRequeue
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	c, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	c.Qos(10, 1, false)
	receiveCh, err := c.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(receiveCh <-chan amqp.Delivery) {
		for delivery := range receiveCh {
			// fmt.Printf("Raw body: %s\n", string(delivery.Body))
			message, err := unmarshaller(delivery.Body)
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

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" 0 or "transient" 1
	handler func(T) AckType, // acktype
) error {
	unmarshaller := func(b []byte) (T, error) {
		var ret T
		err := json.Unmarshal(b, &ret)
		return ret, err
	}
	return subscribe(conn, exchange, queueName, key, simpleQueueType, handler, unmarshaller)
}

func SubscribeGOB[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType, // an enum to represent "durable" 0 or "transient" 1
	handler func(T) AckType, // acktype
) error {
	unmarshaller := func(b []byte) (T, error) {
		decoder := gob.NewDecoder(bytes.NewBuffer(b))
		var ret T
		err := decoder.Decode(&ret)
		return ret, err
	}
	return subscribe(conn, exchange, queueName, key, simpleQueueType, handler, unmarshaller)
}
