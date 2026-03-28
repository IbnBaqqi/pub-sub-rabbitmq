package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON handles publishing message to exchange
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	data, _ := json.Marshal(val)

	ctx := context.Background()
	return ch.PublishWithContext(
		ctx,
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: data,
		})
}
