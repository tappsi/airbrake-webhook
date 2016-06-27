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

func main() {

	port := "8080"
	endpoint := "airbrake-webhook"
	exchange := "notifications_prod"
	uri := "amqp://test:test@192.168.1.13:5672"

	queue := webhook.NewMessagingQueue(uri, exchange)
	hook  := webhook.NewWebHook(queue)

	api := iris.New()
	api.Post("/" + endpoint, hook.Process)
	go cleanup(queue)
	api.Listen(":" + port)

}

func cleanup(queue webhook.MessagingQueue) {

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGTSTP)
	<-sigChan

	fmt.Println("\nReceived an interrupt, stopping services...\n")
	queue.Close()

	runtime.GC()
	os.Exit(0)

}
