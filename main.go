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

	cfg    := webhook.LoadConfiguration("./config/")
	queue  := webhook.NewMessagingQueue(cfg.QueueURI, cfg.ExchangeName, cfg.PoolConfig)
	hook   := webhook.NewWebHook(queue)

	iris.Post("/" + cfg.EndpointName, hook.Process)
	go cleanup(queue)
	iris.Listen(fmt.Sprintf(":%d", cfg.WebServerPort))

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
