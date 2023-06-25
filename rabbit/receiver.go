package rabbit

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"message-service/courier"
	"message-service/helpers"
	"message-service/types"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	channel *amqp.Channel
	queue   string
}

func NewReceiver(ch *amqp.Channel) Receiver {
	receiver := Receiver{
		channel: ch,
	}

	err := declareExchange(ch)
	helpers.CheckForError(err, "Receiver setup failed")

	return receiver
}

func (receiver *Receiver) Consume(db *gorm.DB) {
	queue, err := declareQueue(receiver.channel)
	helpers.CheckForError(err, "Queue declaration failed")

	var forever chan struct{}

	messages, err := receiver.channel.Consume(
		queue.Name,
		"",
		false, // messages should be ack-ed only after success
		false,
		false,
		true, // not forcing other instances to wait
		nil,
	)
	helpers.CheckForError(err, "Failed to register a consumer")

	client := courier.InitEmailClient()
	slackClient := courier.InitSlackClient()

	go func() {
		for d := range messages {
			var content types.Content
			err := json.Unmarshal(d.Body, &content)

			if err != nil {
				// Not re-queueing messages with a broken structure
				// they will be saved in the DB if an inspection is needed
				err = receiver.channel.Nack(d.DeliveryTag, false, false)
				helpers.LogMessageError(err, db, d.Body)

				continue
			}

			switch content.Type {
			case "email":
				err = client.SendMessage(content)
			case "slack":
				time.Sleep(20 * time.Second)
				err = slackClient.SendMessage(content)
			}

			if err != nil {
				helpers.CheckForError(err,
					fmt.Sprintf("Message Sending Failed: Type: %s", content.Type),
				)
				// Requeue failed message
				err = receiver.channel.Nack(d.DeliveryTag, false, true)

				continue
			}

			err = receiver.channel.Ack(d.DeliveryTag, false)

			// An error with acknowledgements means we should drop this instance
			if err != nil {
				log.WithFields(log.Fields{
					"message": "Error with RMQ connection",
				}).Panic(err.Error())
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
