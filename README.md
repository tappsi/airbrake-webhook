This is an Airbrake webhook implemented in Go as a microservice. As part of its integrations, Airbrake allows the definition of a webhook (as described in the [documentation](https://airbrake.io/docs/integrations/webhooks/).) An application using Airbrake is expected to provide an endpoint that gets called by Airbrake whenever a new error occurs, and this implementation provides a generic service to be exposed for that purpose - the received message is forwarded to a queue to be consumed at a later point.

This microservice showcases the usage of several high-performance Go libraries and frameworks (detailed below), and also demonstrates how to send messages to a RabbitMQ queue using a connection pool - the queue's consumer must act on the received message, but this is outside the scope of this project.

# Installation

Assuming that you've installed Go and have correctly set-up the `$GOPATH` and `$GOROOT` environment variables, do the following:

  1. Clone the repository: https://github.com/tappsi/airbrake-webhook.git
  2. Go to the `airbrake-webhook` directory
  3. Install dependencies using `go get -u ./...`
  4. Compile the project `go install github.com/tappsi/airbrake-webhook`

# Usage

Airbrake will invoke the endpoint exposed by this microservice whenever an error occurs. Each invocation will send a new message to our webhook using the following format:

```json
{
  "error":{
    "id":37463546,
    "error_message":"KitchenException: You are all out of bacon!",
    "error_class":"KitchenException",
    "file":"[PROJECT_ROOT]/app/controllers/bacon_controller.rb",
    "line_number":35,
    "project":{
      "id":1111,
      "name":"Baconator"
     },
    "last_notice":{
      "id":4505303522,
      "request_method":null,
      "request_url":"http://airbrake.io:445/bacon/cook",
      "backtrace":[
        "[PROJECT_ROOT]/app/controllers/bacon_controller.rb:35:in `cook'",
        "[PROJECT_ROOT]/app/middleware/kitchen.rb:19:in `oven'",
        "[PROJECT_ROOT]/app/middleware/kitchen.rb:33:in `chef'",
        "[PROJECT_ROOT]/app/middleware/salumi.rb:23:in `store'"
      ]
    },
    "environment":"avocado",
    "first_occurred_at":"2012-02-23T22:03:03Z",
    "last_occurred_at":"2012-03-21T08:37:15Z",
    "times_occurred":118
  },
  "airbrake_error_url": "https://airbrake.io/airbrake-error-url"
}
```

To handle the messages, we have to start the service:

  1. Create the `config/production.json` file with values appropriate for your environment, use `config/development.json` as a guide. In particular, notice that we have to set the endpoint's name, the port where the web server will be running, the RabbitMQ URL and credentials, the exchange to be used and the RabbitMQ connection pool options
  2. Configure Airbrake. Go to your project's Dashboard -> Settings -> Integrations and set the URL where the service is going to live. This is defined based on the parameters configured in the previous step
  3. Go to the `bin` directory of your Go projects and run the executable, passing the appropriate environment. For example: `GO_ENV=production $GOPATH/bin/airbrake-webhook`
  4. The previous step will leave a web server running, listening for calls into the configured endpoint. To kill the service, simply press `Ctrl+C`

I just described a very basic setup for the service, in a real production environment you might have a separate web server acting as proxy and handling secure connections and redirections to the actual service. Also you should take all previsions to ensure that the service never goes down, for instance by using `supervisord` - but this is outside the scope of this document.

Now it's up to you to decide how to handle the message. By default, this microservice parses the JSON with the received message, extracts some of its fields and sends a new message to a RabbitMQ queue for further processing. In our production environment, we have a different process that reads from the queue and sends the message to a log chat for all developers to see. If you want to do something different, modify the `Process()` function in `webhook.go` as required.

# Project Structure



# Dependencies



# Version

Current: **v1.0**

# Author

The author of `airbrake-webhook` is [@oscar-lopez](https://github.com/oscar-lopez), it was developed as an internal project in [Tappsi](https://tappsi.co), and now it's open source.

# Contact

If you have any issues or questions regarding `airbrake-webhook`, feel free to contact the author using any of the following channels, or post a question in StackOverflow:

- [Twitter](https://twitter.com/oscar_lopez)
- [LinkedIn](https://co.linkedin.com/in/óscar-andrés-lópez-436a58)
- [Stackoverflow](http://stackoverflow.com/users/201359/Óscar-lópez)

# License

Unless otherwise noted, the `airbrake-webhook` source files are distributed under the MIT License found in the [LICENSE file](LICENSE).
