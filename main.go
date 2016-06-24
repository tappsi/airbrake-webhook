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

	api := iris.New()
	api.Post("/airbrake-webhook", webhook.Process)

	go cleanup()
	api.Listen(":8080")

}

func cleanup() {

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGTSTP)
	<-sigChan

	fmt.Println("\nReceived an interrupt, stopping services...\n")

	runtime.GC()
	os.Exit(0)

}
