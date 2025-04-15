package workers

import (
	"runtime/debug"

	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
)

func ShutdownOnPanic() {
	if err := recover(); err != nil {
		logger.Errorf("failed to recovered from panic in worker: %v", err)
		logger.Error(string(debug.Stack()))
	}
}
