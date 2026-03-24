package pubsub

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
type SimpleQueueType string

const (
	DurableQueue SimpleQueueType = "durable"
	TransientQueue SimpleQueueType = "transient"
)

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
		return nil, amqp.Queue{}, err
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
		log.Printf("unable to declare queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	err = amqpCh.QueueBind(
		queue.Name,  // queue name
		key,         // routing key
		exchange,    // exchange
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Printf("unable to bind queue: %v", err)
		return nil, amqp.Queue{}, err
	}

	return amqpCh, queue, nil
}