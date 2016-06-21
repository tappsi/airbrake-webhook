package webhook

import (
	"bytes"
	"strconv"
	"strings"
	"github.com/kataras/iris"
	"github.com/buger/jsonparser"
	"github.com/mailru/easyjson/jwriter"
)

// easyjson:json
type Notification struct {
	service    string
	recipients []string
	message    string
}

func Process(ctx *iris.Context) {

	input := ctx.GetRequestCtx().Request.Body()

	// parse input

	environment, _, _, _ := jsonparser.Get(input, "error", "environment")
	timesOccurred, _ := jsonparser.GetInt(input, "error", "times_occurred")
	errorId, _ := jsonparser.GetInt(input, "error", "last_notice", "id")
	errorUrl, _, _, _ := jsonparser.Get(input, "airbrake_error_url")
	errorMessage, _, _, _ := jsonparser.Get(input, "error", "error_message")

	// create output

	service := "mattermost"
	recipients := []string{"opslog"}
	message := strings.Join(
		[]string{
			"Environment: " + string(environment),
			"Occurrences: " + strconv.FormatInt(timesOccurred, 10),
			"Error ID: " + strconv.FormatInt(errorId, 10),
			"Error URL: " + string(errorUrl),
			"Error Message: " + string(errorMessage),
		},
		", ")

	notification := Notification{service, recipients, message}
	writer := jwriter.Writer{}
	notification.MarshalEasyJSON(&writer)

	if writer.Error == nil {
		buf := new(bytes.Buffer)
		writer.DumpTo(buf)
		SendMessage(buf.Bytes())
	}

}
