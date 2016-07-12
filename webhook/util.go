package webhook

import (
	"fmt"
	"log"
	"strings"
)

func IsEmpty(txt string) bool {
	return strings.TrimSpace(txt) == ""
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
