package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/rs/zerolog"
)

type attachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type channelNotification struct {
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Channel     string       `json:"channel"`
	Attachments []attachment `json:"attachments"`
}

type slackHook struct {
	client        *http.Client
	SlackURL      string
	SlackUsername string
	wg            *sync.WaitGroup
}

func NewSlackHook(wg *sync.WaitGroup) *slackHook {
	return &slackHook{
		client:        &http.Client{},
		SlackURL:      os.Getenv("SLACK_URL"),
		SlackUsername: os.Getenv("SLACK_USERNAME"),
		wg:            wg,
	}
}

func (p *slackHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	log.Printf("found message for sending: [%+v]", level)

	if level >= zerolog.ErrorLevel {
		notification := &channelNotification{
			Text:      "Error Notification", // Move this to appErr and fetching it here after type casting
			Username:  p.SlackUsername,
			IconEmoji: ":error:",
			Channel:   "errors",
			Attachments: []attachment{
				{
					Title: "Error Notification",
					Text:  message,
				},
			},
		}

		log.Printf("sending informatics to telegram [%+v]", true)

		err := p.sendChannelNotification(notification)
		if err != nil {
			log.Printf("failed to send a message: [%+v]", err.Error())
		}
	}
}

func (p *slackHook) sendChannelNotification(
	notification *channelNotification,
) error {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return apperrs.Wrap(
			err,
		).AddLogMessage(
			"marshall slack notification",
		)
	}

	request, err := http.NewRequest(
		"POST",
		p.SlackURL,
		bytes.NewBuffer(notificationBytes),
	)
	if err != nil {
		return apperrs.Wrap(
			err,
		).AddLogMessage(
			"create slack request",
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// response, err := p.client.Do(request)
	// if err != nil {
	// 	return apperrs.Wrap(
	// 		err,
	// 	).AddLogMessage(
	// 		"make slack notification request",
	// 	)
	// }

	// if response.StatusCode < 200 && response.StatusCode > 299 {
	// 	return apperrs.Wrap(
	// 		errors.New("unable to send notification"),
	// 	).AddLogMessage(
	// 		"slack request error",
	// 	)
	// }

	return nil
}
