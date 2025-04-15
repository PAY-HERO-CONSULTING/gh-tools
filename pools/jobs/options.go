package jobs

import "github.com/gocraft/work"

type JobOption func(work.JobOptions)

func WithHighPriority() JobOption {
	return WithPriority(highPriority)
}

func WithLowPriority() JobOption {
	return WithPriority(lowPriority)
}

func WithMaxConcurrency(maxConcurrency uint) JobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.MaxConcurrency = maxConcurrency
	}
}

func WithMaxFails(maxFails uint) JobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.MaxFails = maxFails
	}
}

func WithPriority(priority uint) JobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.Priority = priority
	}
}

func WithSkipDeadQueue(skip bool) JobOption {
	return func(jobOptions work.JobOptions) {
		jobOptions.SkipDead = skip
	}
}
