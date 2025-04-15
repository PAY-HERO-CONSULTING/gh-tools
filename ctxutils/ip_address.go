package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
)

func IPAddress(ctx context.Context) string {
	existing := ctx.Value(dtos.ContextKeyIpAddress)
	if existing == nil {
		return ""
	}

	return existing.(string)
}

func WithIpAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, dtos.ContextKeyIpAddress, ipAddress)
}
