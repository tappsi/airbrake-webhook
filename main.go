package main

import "github.com/kataras/iris"
import "github.com/tappsi/airbrake-webhook/webhook"

func main() {
  api := iris.New()
  api.Post("/airbrake", webhook.Process)
  api.Listen(":8080")
}
