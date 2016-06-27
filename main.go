package main

import (
	"os"
	"fmt"
	"runtime"
	"syscall"
	"os/signal"
	"github.com/kataras/iris"
	"github.com/tappsi/airbrake-webhook/webhook"
)

var connPool webhook.RMQConnectionPool

func main() {

	api := iris.New()
	api.Post("/airbrake-webhook", webhook.Process)

	uri  := "amqp://test:test@192.168.1.13:5672"
	connPool = webhook.NewRMQConnectionPool(uri)

	conn, toReturn, _ := connPool.GetConnection()
	conn.Close()
	_ := connPool.ReturnConnection(toReturn)

	go cleanup()
	api.Listen(":8080")

}

func cleanup() {

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGTSTP)
	<-sigChan

	fmt.Println("\nReceived an interrupt, stopping services...\n")
	connPool.ClosePool()

	runtime.GC()
	os.Exit(0)

}
