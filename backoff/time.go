package backoff

import (
	"math/rand"
	"time"
)

var commonBackoffSequence = []int{300, 420, 660, 780, 1020, 1140, 1380, 1740, 1860, 2220}
var fasterBackoffSequence = []int{5, 7, 13, 23, 37}

func backoffRetryDelay(retry int, backoffSequence []int) time.Duration {

	duration := backOff(retry, backoffSequence)

	rand.Seed(time.Now().UnixNano())
	durationWithJitter := rand.Intn(duration) + duration

	return time.Duration(durationWithJitter) * time.Second
}

func backOff(retry int, backoffSequence []int) int {

	increment := retry % len(backoffSequence)

	if increment > len(backoffSequence)-1 {
		increment = 0
	}

	factor := backoffSequence[increment]

	return factor
}

func CommomBackoffDelay(retry int) time.Duration {
	return backoffRetryDelay(retry, commonBackoffSequence)
}

func FasterBackoffDelay(retry int) time.Duration {
	return backoffRetryDelay(retry, fasterBackoffSequence)
}
