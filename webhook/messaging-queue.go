package webhook

import (
	"github.com/streadway/amqp"
	"log"
)

// MessagingQueue is a structure used for defining a messaging
// queue object, it's used for encapsulating the exchange
// name and the connection pool used for connecting to RMQ.
type MessagingQueue struct {
	pool     *RMQConnectionPool
	exchange string
}

// NewMessagingQueue creates a new MessagingQueue object given as parameters
// the uri where it's located (including authentication credentials), the
// name of the exchange and a PoolConfiguration struct with the pool's config.
func NewMessagingQueue(uri, exchange string, cfg PoolConfiguration) MessagingQueue {
	pool := NewRMQConnectionPool(uri, cfg)
	return MessagingQueue{pool: &pool, exchange: exchange}
}

// SendMessage sends a message to RMQ. The message body is received as
// parameter. Internally, this method obtains a connection from the pool,
// marshals the message body into a suitable JSON and publishes it to
// the specified RMQ exchange, with QoS considerations. In particular,
// a confirmation channel is defined for receiving message acknowledges.
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

// Close closes the queue's connections still open
func (m *MessagingQueue) Close() {
	m.pool.Close()
}

// freeResources is a deferred method for freeing resources after sending
// a message to the queue, basically it returns a connection to the pool.
func (m *MessagingQueue) freeResources(toReturn interface{}, ch *amqp.Channel) {
	e1 := ch.Close()
	FailOnError(e1, "Error closing channel")
	e2 := m.pool.ReturnConnection(toReturn)
	FailOnError(e2, "Error returning connection to pool")
}

// confirmOne is a deferred function for receiving message acknowledges from the queue.
func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; !confirmed.Ack {
		log.Fatalf("Failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
