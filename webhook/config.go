package webhook

import (
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
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

	file := getFile(path)
	raw, err := ioutil.ReadFile(file)
	FailOnError(err, "Can not load configuration")

	var cfg Configuration
	err = json.Unmarshal(raw, &cfg)
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
