package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
type SimpleQueueType string

const (
	DurableQueue SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

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

// DeclareAndBind create a channel on conn and binds the exchange to a queue
// with routing key
func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {

	// create a channel on the connection
	amqpCh, err := conn.Channel()
	if err != nil {
		log.Printf("unable to create channel: %v", err)
	}

	durable := false
	if queueType == DurableQueue {
		durable = true
	}
	autodelete, exclusive := false, false
	if queueType == TransientQueue {
		autodelete = true
		exclusive = true
	}

	// declear a new queue with it properties to hold messages & deliver to consumer
	queue, err := amqpCh.QueueDeclare(
		queueName,
		durable,
		autodelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		log.Printf("unable to declare queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	err = amqpCh.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Printf("unable to bind queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	return amqpCh, queue, nil
}
