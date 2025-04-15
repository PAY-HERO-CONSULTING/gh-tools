package dtos

import "strings"

type Notifications struct {
	Single bool   `json:"single"`
	Type   string `json:"type"`
	Data   any    `json:"data"`
}

func (n *Notifications) IsMail() bool {
	return strings.ToLower(n.Type) == "mail"
}

func (n *Notifications) IsSMS() bool {
	return strings.ToLower(n.Type) == "sms"
}

func (n *Notifications) Event() bool {
	return strings.ToLower(n.Type) == "event"
}
