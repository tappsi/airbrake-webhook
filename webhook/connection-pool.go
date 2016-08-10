package webhook

import (
	pool "github.com/jolestar/go-commons-pool"
	"github.com/streadway/amqp"
)

// connection pool

// RMQConnectionPool is the struct used for defining
// an RMQ connection pool object.
type RMQConnectionPool struct {
	pool *pool.ObjectPool
}

// NewRMQConnectionPool creates a new connection pool object
// for connecting with RMQ. It receives as parameters the
// uri used for accessing RMQ (including authentication
// credentials) and cfg, the pool configuration.
func NewRMQConnectionPool(uri string, cfg PoolConfiguration) RMQConnectionPool {
	pool := pool.NewObjectPoolWithDefaultConfig(&RMQConnectionFactory{uri: uri})
	pool.Config.Lifo = false
	pool.Config.MaxTotal = cfg.MaxTotal
	pool.Config.MinIdle = cfg.MinIdle
	pool.Config.MaxIdle = cfg.MaxIdle
	return RMQConnectionPool{pool: pool}
}

// GetConnection obtains a new amqp connection from the pool.
func (p *RMQConnectionPool) GetConnection() (*amqp.Connection, interface{}, error) {
	obj, err := p.pool.BorrowObject()
	conn := obj.(**amqp.Connection)
	return *conn, obj, err
}

// ReturnConnection returns an amqp connection to the pool. It receives
// as parameter the connection to be returned to the pool.
func (p *RMQConnectionPool) ReturnConnection(conn interface{}) error {
	return p.pool.ReturnObject(conn)
}

// Close closes the amqp connection pool object, freeing its resources.
func (p *RMQConnectionPool) Close() {
	p.pool.Close()
}

// connection factory

// RMQConnectionFactory is the struct used for defining
// an RMQ connection factory object. The methods of this
// struct implement the PooledObject interface.
type RMQConnectionFactory struct {
	uri string
}


// MakeObject creates a new amqp connection and encapsulates it as
// a pooled object. This method is part of the PooledObject interface.
func (f *RMQConnectionFactory) MakeObject() (*pool.PooledObject, error) {
	conn, err := amqp.Dial(f.uri)
	return pool.NewPooledObject(&conn), err
}

// MakeObject closes an amqp connection and encapsulates it as
// a pooled object. It receives as parameter a PooledObject.
// This method is part of the PooledObject interface.
func (f *RMQConnectionFactory) DestroyObject(object *pool.PooledObject) error {
	conn := object.Object.(**amqp.Connection)
	return (*conn).Close()
}

// ValidateObject is part of the PooledObject interface, currently is not
// being used. It receives as parameter a PooledObject.
func (f *RMQConnectionFactory) ValidateObject(object *pool.PooledObject) bool {
	return true
}

// ActivateObject is part of the PooledObject interface, currently is not
// being used. It receives as parameter a PooledObject.
func (f *RMQConnectionFactory) ActivateObject(object *pool.PooledObject) error {
	return nil
}

// PassivateObject is part of the PooledObject interface, currently is not
// being used. It receives as parameter a PooledObject.
func (f *RMQConnectionFactory) PassivateObject(object *pool.PooledObject) error {
	return nil
}
