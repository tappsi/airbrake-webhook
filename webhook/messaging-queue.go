package webhook

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

func SendMessage(body []byte) bool {

	uri := "amqp://test:test@192.168.1.13:5672"
	exchange := "notifications_test"

	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a exchange")

	err = ch.Confirm(false)
	failOnError(err, "Channel could not be put into confirm mode")
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer confirmOne(confirms)

	err = ch.Publish(
		exchange, // publish to an exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
	failOnError(err, "Failed to publish message")

	return true

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; !confirmed.Ack {
		log.Fatalf("Failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
