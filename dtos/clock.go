package dtos

import "time"

type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type realClock struct{}

func NewClock() Clock {
	return &realClock{}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Since(
	t time.Time,
) time.Duration {
	return time.Since(t)
}
