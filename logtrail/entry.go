package logtrail

import (
	"time"
)

// Entry holds all information related to an API call event.
type Entry struct {
	Agent      string    `json:"agent"`
	UserID     string    `json:"user_id"`
	Time       time.Time `json:"time"`
	RequestID  string    `json:"request_id"`
	MacAddress string    `json:"mac_address"`
	HTTP       HTTPEntry `json:"http_entry"`
}

// HTTPEntry contains details related to an HTTP call for an audit log entry.
type HTTPEntry struct {
	ClientIP     string   `json:"client_ip"`
	UserAgent    string   `json:"user_agent"`
	Method       string   `json:"method"`
	Path         string   `json:"path"`
	RequestBody  string   `json:"request_body"`
	StatusCode   int      `json:"status_code"`
	ResponseTime int      `json:"response_time"`
	ResponseSize int      `json:"response_size"`
	Errors       []string `json:"errors"`
}
