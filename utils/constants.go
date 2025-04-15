package utils

import (
	"math/rand"
	"time"
)

const (
	digits                         = "0123456789"
	hexLetters                     = "abcdef"
	lowerCaseLetters               = "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetters               = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	africasTalkingInboxTimeLayout  = "2006-01-02T15:04:05Z"
	africasTalkingInboxTimeLayout2 = "2006-01-02 15:04:05"
	customDateTimeLayout           = "02-01-2006 15:04:05"
	dateLayout                     = "2006-01-02"
	mpesaTimeLayout                = "20060102150405"
	paypalTimeLayout               = "15:04:05 Jan 02, 2006 MST"
	safaricomTimeLayout            = "2006-01-02 15:04:05"
	timeLayout                     = "2006/01/02 15:04:05"
	timeandtimezonelayout          = "2006/01/02 15:04:05-07:00"
)

var generalBackoffSequence = []int{300, 420, 660, 780, 1020, 1140, 1380, 1740, 1860, 2220}
var fasterBackoffSequence = []int{5, 7, 13, 23, 37}
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
