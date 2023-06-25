package rabbit

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		os.Getenv("RMQ_EXCHANGE"),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		os.Getenv("RMQ_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
}
