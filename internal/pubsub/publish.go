package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)	
	err := encoder.Encode(val)
	if err != nil {
		return err
	}
	
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	})
}

func PublishGameLog(gl routing.GameLog, ch *amqp.Channel) error {
	return PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug + "." + gl.Username, gl)
}

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	return nil
}
