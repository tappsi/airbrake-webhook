package webhook

import (
	"encoding/json"
	"os"
	"strings"
)

type Configuration struct {
	WebServerPort uint16            `json:"webserver-port"`
	EndpointName  string            `json:"endpoint-name"`
	ExchangeName  string            `json:"exchange-name"`
	QueueURI      string            `json:"queue-uri"`
	PoolConfig    PoolConfiguration `json:"pool-config"`
}

type PoolConfiguration struct {
	MaxTotal int `json:"max-total"`
	MinIdle  int `json:"min-idle"`
	MaxIdle  int `json:"max-idle"`
}

func LoadConfiguration(path string) Configuration {

	fileName := getFile(path)
	file, err := os.Open(fileName)
	FailOnError(err, "Can not load configuration")

	var cfg Configuration
	err = json.NewDecoder(file).Decode(&cfg)
	FailOnError(err, "Can not parse configuration")

	return cfg

}

func getFile(path string) string {

	env := os.Getenv("GO_ENV")

	if IsEmpty(env) {
		env = "development"
	} else {
		env = strings.ToLower(strings.TrimSpace(env))
	}

	return path + env + ".json"

}
