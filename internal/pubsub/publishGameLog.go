package pubsub

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(gl routing.GameLog, ch *amqp.Channel) error {
	return PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug + "." + gl.Username, gl)
}
