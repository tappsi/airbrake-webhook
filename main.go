package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/tappsi/airbrake-webhook/webhook"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	cfg := webhook.LoadConfiguration("./config/")
	queue := webhook.NewMessagingQueue(cfg.QueueURI, cfg.ExchangeName, cfg.PoolConfig)
	hook := webhook.NewWebHook(queue)

	iris.Post("/" + cfg.EndpointName, hook.Process)
	go cleanup(queue)

	err := iris.ListenTo(config.Server{ListeningAddr: fmt.Sprintf(":%d", cfg.WebServerPort)})
	webhook.FailOnError(err, "Error listening on web server")

}

func cleanup(queue webhook.MessagingQueue) {

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGTSTP)
	<-sigChan

	fmt.Println("\nReceived an interrupt, stopping services...\n")
	queue.Close()
	iris.Close()

	runtime.GC()
	os.Exit(0)

}
