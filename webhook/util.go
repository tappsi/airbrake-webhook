package webhook

import (
	"fmt"
	"log"
	"strings"
)

// IsEmpty returns a boolean indicating whether
// the txt parameter is an empty string.
func IsEmpty(txt string) bool {
	return strings.TrimSpace(txt) == ""
}

// FailOnError logs an error if the given parameter
// err is not null, using the msg message parameter.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
