package ctxutils

import (
	"context"
	"time"
)

func ExtendContext(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {

	ctx, cancel := context.WithTimeout(ctx, duration)

	return ctx, cancel
}
