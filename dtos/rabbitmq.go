package dtos

import "time"

type Delivery struct {
	Body            []byte
	ContentEncoding string
	ContentType     string
	Exchange        string
	Expiration      string
	Headers         map[string]interface{}
	MessageId       string
	RoutingKey      string
	Timestamp       time.Time
	Type            string
}
