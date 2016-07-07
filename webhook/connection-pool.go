package webhook

import (
	"github.com/streadway/amqp"
	pool "github.com/jolestar/go-commons-pool"
)

// connection pool

type RMQConnectionPool struct {
	pool *pool.ObjectPool
}

func NewRMQConnectionPool(uri string, cfg PoolConfiguration) RMQConnectionPool {
	pool := pool.NewObjectPoolWithDefaultConfig(&RMQConnectionFactory{uri: uri})
	pool.Config.Lifo     = false
	pool.Config.MaxTotal = cfg.MaxTotal
	pool.Config.MinIdle  = cfg.MinIdle
	pool.Config.MaxIdle  = cfg.MaxIdle
	return RMQConnectionPool{ pool: pool }
}

func (p *RMQConnectionPool) GetConnection() (*amqp.Connection, interface{}, error) {
	obj, err := p.pool.BorrowObject()
	conn := obj.(**amqp.Connection)
	return *conn, obj, err
}

func (p *RMQConnectionPool) ReturnConnection(conn interface{}) error {
	return p.pool.ReturnObject(conn)
}

func (p *RMQConnectionPool) Close() {
	p.pool.Close()
}

// connection factory

type RMQConnectionFactory struct {
	uri string
}

func (f *RMQConnectionFactory) MakeObject() (*pool.PooledObject, error) {
	conn, err := amqp.Dial(f.uri)
	return pool.NewPooledObject(&conn), err
}

func (f *RMQConnectionFactory) DestroyObject(object *pool.PooledObject) error {
	conn := object.Object.(**amqp.Connection)
	return (*conn).Close()
}

func (f *RMQConnectionFactory) ValidateObject(object *pool.PooledObject) bool {
	return true
}

func (f *RMQConnectionFactory) ActivateObject(object *pool.PooledObject) error {
	return nil
}

func (f *RMQConnectionFactory) PassivateObject(object *pool.PooledObject) error {
	return nil
}
