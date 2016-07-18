package webhook

import (
	"bytes"
	"github.com/buger/jsonparser"
	"github.com/kataras/iris"
	"github.com/mailru/easyjson/jwriter"
	"strconv"
	"strings"
)

// easyjson:json
type Notification struct {
	service    string
	recipients []string
	message    string
}

type WebHook struct {
	queue *MessagingQueue
}

func NewWebHook(queue MessagingQueue) WebHook {
	return WebHook{queue: &queue}
}

func (w *WebHook) Process(ctx *iris.Context) {

	input := ctx.Request.Body()

	// parse input

	environment, _, _, _ := jsonparser.Get(input, "error", "environment")
	timesOccurred, _ := jsonparser.GetInt(input, "error", "times_occurred")
	errorId, _ := jsonparser.GetInt(input, "error", "id")
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
			"Error URL: " + strings.Replace(string(errorUrl), "\\", "", -1),
			"Error Message: " + string(errorMessage),
		},
		", ")

	notification := Notification{service, recipients, message}
	writer := jwriter.Writer{}
	notification.MarshalEasyJSON(&writer)

	if writer.Error == nil {
		buf := new(bytes.Buffer)
		writer.DumpTo(buf)
		w.queue.SendMessage(buf.Bytes())
	}

}
