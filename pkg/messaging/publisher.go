package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func PublishMessageJSON(queue string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return PublishMessage(queue, string(data))
}

func PublishMessage(queue string, message string) error {
	conn, err := amqp091.Dial(generateDSN())
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	if err := ch.Confirm(false); err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		return err
	}
	select {
	case ntf := <-confirms:
		if !ntf.Ack {
			return fmt.Errorf("message delivery failed")
		}
	case <-ch.NotifyReturn(make(chan amqp091.Return)):
		return fmt.Errorf("message delivery failed")
	}

	return nil
}
