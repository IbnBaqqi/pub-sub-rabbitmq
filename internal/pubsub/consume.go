package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
type SimpleQueueType string

const (
	DurableQueue SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

// DeclareAndBind opens a new channel on conn, declares a queue with the given
// queueType (durable or transient), and binds it to exchange using key as the
// routing key. Returns the channel and declared queue.
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
		return nil, amqp.Queue{}, fmt.Errorf("unable to create channel: %v", err)
	}

	// declear a new queue with it properties to hold messages & deliver to consumer
	queue, err := amqpCh.QueueDeclare(
		queueName,                     // name
		queueType == DurableQueue,     // durable
		queueType == TransientQueue,   // delete when unused
		queueType == TransientQueue,   // exclusive
		false,                         // no-wait
		nil,                           // args
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("unable declare queue: %v", err)
	}

	err = amqpCh.QueueBind(
		queue.Name,  // queue name
		key,         // routing key
		exchange,    // exchange
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("unable to bind queue: %v", err)
	}

	return amqpCh, queue, nil
}

// SubscribeJSON declares and binds a queue, then consumes messages from it,
// unmarshalling each delivery as JSON into type T and passing it to handler.
// Blocks until the channel is closed.
func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType,
    handler func(T),
) error {

	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveriesChan, err := ch.Consume(queueName, "", false, false, false, false, nil)

	var value T
	for delivery := range deliveriesChan {
		if err := json.Unmarshal(delivery.Body, &value); err != nil {
			return err
		}
		handler(value)
		if err := delivery.Ack(false); err != nil {
			return err 
		}
	}

	return nil
}