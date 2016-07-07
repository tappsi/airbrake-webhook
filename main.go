package main

import (
	"os"
	"fmt"
	"runtime"
	"syscall"
	"os/signal"
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/secure"
	"github.com/tappsi/airbrake-webhook/webhook"
)

func main() {

	cfg    := webhook.LoadConfiguration("./config/")
	queue  := webhook.NewMessagingQueue(cfg.QueueURI, cfg.ExchangeName, cfg.PoolConfig)
	hook   := webhook.NewWebHook(queue)
	secure := secureConfig(cfg.SecureConfig)

	iris.UseFunc(func(c *iris.Context) {
		err := secure.Process(c)
		webhook.FailOnError(err, "Failed to config secure middleware")
		c.Next()
	})
	iris.Post("/" + cfg.EndpointName, hook.Process)

	go cleanup(queue)
	iris.Listen(fmt.Sprintf(":%d", cfg.WebServerPort))

}

func secureConfig(cfg webhook.SecureConfiguration) *secure.Secure {
	return secure.New(secure.Options{
		IsDevelopment: cfg.IsDevelopment,
	})
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
