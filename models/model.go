package models

import (
	"time"
)

type TokenInfo struct {
	AccountIDs     []string
	Exp            time.Time
	IsAdmin        bool
	OrganizationID string
	Refresh        time.Time
	Status         string
	UserID         string
	IsWorker       bool
	SessionID      string
	Roles          map[string]string
	IsLoggedIn     bool
	Token          string
	IsAPIkey       bool
	IsPinSet       bool
	PhoneVerified  bool
	EmailVefifed   bool
	Scope          string
	Username       string
}
