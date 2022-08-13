package messaging

import (
	"github.com/rabbitmq/amqp091-go"
)

func BuildConsumer(queue string, handler func(string) (bool, error)) error {
	conn, err := amqp091.Dial(generateDSN())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	defer func() {
		_ = conn.Close()
		_ = ch.Close()
	}()

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			requeue, err := handler(string(d.Body))
			if err != nil {
				_ = d.Nack(false, requeue)
			} else {
				_ = d.Ack(false)
			}
		}
	}()
	return nil
}
