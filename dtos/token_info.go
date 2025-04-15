package dtos

import (
	"time"
)

type TokenInfo struct {
	AccountIDs     []int64
	Exp            time.Time
	IsAdmin        bool
	OrganizationID int64
	Refresh        time.Time
	Status         string
	UserID         int64
	IsWorker       bool
	JID            string
	Roles          map[int64]string
}
