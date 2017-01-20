# Airbrake Webhook

This is an Airbrake webhook implemented in Go as a microservice. As part of its integrations, Airbrake allows the definition of a webhook (as described in the [documentation](https://airbrake.io/docs/integrations/webhooks/).) An application using Airbrake is expected to provide an endpoint that gets called by Airbrake whenever a new error occurs, and this implementation provides a generic service to be exposed for that purpose - the received message is forwarded to a messaging exchange to be consumed asynchronously at a later point.

This microservice showcases the usage of several high-performance Go libraries and frameworks (detailed below), and also demonstrates how to send messages to a RabbitMQ exchange using a connection pool - the exchange's asynchronous consumer(s) must act on the received message, but that's outside the scope of this project.

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

  1. Create the `config/production.json` file with values appropriate for your environment, use `config/development.json` as a guide, see the "Configuration" section below for more details.
  2. Configure Airbrake. Go to your project's Dashboard -> Settings -> Integrations and set the URL where the service is going to live. This is defined based on the parameters configured in the previous step.
  3. Go to the `bin` directory of your Go projects and run the executable, passing the appropriate environment. For example: `GO_ENV=production $GOPATH/bin/airbrake-webhook`
  4. The previous step will leave a web server running, listening for calls into the configured endpoint. To kill the service, simply press `Ctrl+C`.

I just described a very basic setup for the service, in a real production environment you might have a separate web server acting as proxy and handling secure connections and redirections to the actual service. Also you should take all previsions to ensure that the service never goes down, for instance by using `supervisord` - but this is outside the scope of this document.

Now it's up to you to decide how to handle the message. By default, this microservice parses the JSON with the received message, extracts some of its fields and sends a new message to a RabbitMQ exchange for further processing; the newly created message has the following structure:

```json
{
  "service": "mattermost",
  "recipients": ["opslog"],
  "message": "Environment: production,
              Occurrences: 1,
              Error ID: 37463546,
              Error URL: https://airbrake.io/airbrake-error-url,
              Error Message: missing attribute"
}
```

In our production environment, we have a different process that reads from the exchange and sends the message to a Mattermost log chat for all developers to see. If you want to do something different, modify the `Process()` function in `webhook.go` as required.

# Project Structure

The project is structured as follows:

- A `main.go` file in the `main` package.
- All the `.go` logic files reside inside the `webhook` package.
- Configuration files in `.json` format, residing in the `config` directory.

## Main Package

The `main` package contains a single file, `main.go`, in charge of starting up the application:

- `main.go`: Creates a pool of RMQ connections, starts the web server in charge of handling requests to the exposed endpoint and closes the pool when the application exits.

## Webhook Package

The `main` package contains all the files implementing the application's logic:

- `config.go`: Reads and parses JSON configuration files using Go's standard JSON libraries. The appropriate file is processed depending on the execution environment set in the `GO_ENV` environment variable. If `GO_ENV` was not specified, uses the `"development"` environment by default.
- `connection-pool.go`: Adapts a generic pool library for handling a pool of RMQ's connections. Implements the `PooledObject` interface, required by `go-commons-pool`.
- `easyjson.go`: File auto-generated using the `easyjson` library, customized for performing fast encoding/decoding of notification JSONs.
- `messaging-queue.go`: Implements the functionality for creating a connection to a RabbitMQ server, sending messages to it and freeing-up resources afterwards. The connections are handled using a pool; see `connection-pool.go`. Several QoS attributes can be configured, check RMQ's and `amqp`'s documentation.
- `util.go`: Utility functions for dealing with strings, errors, etc.
- `webhook.go`: The actual webhook. Implements an `iris` handler for processing an http request, parses the JSON message sent by Airbrake, creates a new JSON message with the relevant fields and publishes it to a RabbitMQ exchange for further asynchronous processing.

## Configuration

The webhook's configuration is handled via JSON files, located in the `config` directory. Create one for each of the environments, for instance: `development.json`, `production.json`. The expected format is:

```json
{
  "webserver-port": 8181,
  "endpoint-name": "airbrake-webhook",
  "exchange-name": "notifications_dev",
  "queue-uri": "amqp://test:test@192.168.1.13:5672",
  "pool-config": {
    "max-total": 10,
    "min-idle":   0,
    "max-idle":  10
  }
}
```

Each of the configuration options is detailed below:

* `webserver-port`: Port where the web server runs.
* `endpoint-name`: Path name of the exposed service.
* `exchange-name`: Name of the exchange used to publish messages.
* `queue-uri`: URL and credentials for RMQ.
* `max-total`: Maximum number of total connections open.
* `min-idle`: Minimum number of idle connections allowed.
* `max-idle`: Maximum number of idle connections allowed.

# Dependencies

Several high-performance, third party open source libraries and frameworks were used for writing this project. This was done for showcasing the best libraries for each job, given Go's emphasis on high-speed processing:

* [`iris`](https://github.com/kataras/iris): According to Iris' benchmarks, this is the fastest web framework for Go, starting up its own web server. It's used for exposing a RESTful endpoint consumed by Airbrake. Bear in mind that Iris is a recent framework and somewhat unstable.
* [`amqp`](https://github.com/streadway/amqp): The standard Go client for AMQP, used for connecting to a RabbitMQ server and sending messages to it.
* [`go-commons-pool`](https://github.com/jolestar/go-commons-pool): A generic object pool for Go, the connection to RabbitMQ are pooled using this library.
* [`jsonparser`](https://github.com/buger/jsonparser): Alternative JSON parser for Go that does not require schema, it's the fastest parser for decoding a JSON object. Used for decoding the message sent by Airbrake.
* [`easyjson`](https://github.com/mailru/easyjson): Fast JSON serializer for Go, used for encoding the message sent to a RMQ exchange.

# Version

Current: **v1.0.0**

# Author

The author of `airbrake-webhook` is [@oscar-lopez](https://github.com/oscar-lopez), it was developed as an internal project in [Tappsi](https://tappsi.co), and now it's open source.

# Contact

If you have any issues or questions regarding `airbrake-webhook`, feel free to contact the author using any of the following channels, or post a question in StackOverflow:

- [Twitter](https://twitter.com/oscar_lopez)
- [LinkedIn](https://co.linkedin.com/in/óscar-andrés-lópez-436a58)
- [StackOverflow](http://stackoverflow.com/users/201359/Óscar-lópez)

# License

Unless otherwise noted, the `airbrake-webhook` source files are distributed under the MIT License found in the [LICENSE file](LICENSE).
