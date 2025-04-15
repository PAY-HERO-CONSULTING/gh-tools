package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
)

func UserAgent(ctx context.Context) string {
	existing := ctx.Value(dtos.ContextKeyUserAgent)
	if existing == nil {
		return ""
	}

	return existing.(string)
}

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, dtos.ContextKeyUserAgent, userAgent)
}
