package webhook

import (
	"log"
	"github.com/streadway/amqp"
)

type MessagingQueue struct {
	pool *RMQConnectionPool
	exchange string
}

func NewMessagingQueue(uri, exchange string, cfg PoolConfiguration) MessagingQueue {
	pool := NewRMQConnectionPool(uri, cfg)
	return MessagingQueue{ pool: &pool, exchange: exchange }
}

func (m *MessagingQueue) SendMessage(body []byte) bool {

	conn, obj, err := m.pool.GetConnection()
	FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer m.freeResources(obj, ch)

	err = ch.ExchangeDeclare(
		m.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	)
	FailOnError(err, "Failed to declare a exchange")

	err = ch.Confirm(false)
	FailOnError(err, "Channel could not be put into confirm mode")
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer confirmOne(confirms)

	err = ch.Publish(
		m.exchange, // publish to an exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
	FailOnError(err, "Failed to publish message")

	return true

}

func (m *MessagingQueue) Close() {
	m.pool.Close()
}

func (m *MessagingQueue) freeResources(toReturn interface{}, ch *amqp.Channel) {
	e1 := ch.Close()
	FailOnError(e1, "Error closing channel")
	e2 := m.pool.ReturnConnection(toReturn)
	FailOnError(e2, "Error returning connection to pool")
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; !confirmed.Ack {
		log.Fatalf("Failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
