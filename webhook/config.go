package webhook

import (
	"encoding/json"
	"os"
	"strings"
)

// Configuration is the struct used for parsing the
// JSON used for general configuration
type Configuration struct {
	WebServerPort uint16            `json:"webserver-port"`
	EndpointName  string            `json:"endpoint-name"`
	ExchangeName  string            `json:"exchange-name"`
	QueueURI      string            `json:"queue-uri"`
	PoolConfig    PoolConfiguration `json:"pool-config"`
}

// PoolConfiguration is the struct used for parsing the
// JSON used for connection pool configuration
type PoolConfiguration struct {
	MaxTotal int `json:"max-total"`
	MinIdle  int `json:"min-idle"`
	MaxIdle  int `json:"max-idle"`
}

// LoadConfiguration parses the service's configuration from
// a JSON file, a different file name is used depending on
// the value of the GO_ENV environment variable - for example:
// "production", "development". The path parameter indicates
// where the file is located.
func LoadConfiguration(path string) Configuration {

	fileName := getFile(path)
	file, err := os.Open(fileName)
	FailOnError(err, "Can not load configuration")

	var cfg Configuration
	err = json.NewDecoder(file).Decode(&cfg)
	FailOnError(err, "Can not parse configuration")

	return cfg

}

// getFile obtains the name of the file to load depending on
// the value set on the GO_ENV environment variable. If
// none was set, "development" is used as default.
// The file path is passed as parameter.
func getFile(path string) string {

	env := os.Getenv("GO_ENV")

	if IsEmpty(env) {
		env = "development"
	} else {
		env = strings.ToLower(strings.TrimSpace(env))
	}

	return path + env + ".json"

}
