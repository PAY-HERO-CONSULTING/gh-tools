package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
)

func AllowLeadingWhitespace(ctx context.Context) bool {
	allowLeadingWhitespace := ctx.Value(dtos.ContextKeyLeadingWhitespace)
	if allowLeadingWhitespace == nil {
		return false
	}

	val, ok := allowLeadingWhitespace.(bool)
	if !ok {
		return false
	}

	return val
}

func WithLeadingWhitespace(ctx context.Context, allowLeadingWhitespace bool) context.Context {
	return context.WithValue(ctx, dtos.ContextKeyLeadingWhitespace, allowLeadingWhitespace)
}
