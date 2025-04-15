package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
)

func AccountID(ctx context.Context) int64 {
	existing := ctx.Value(dtos.ContextKeyAccountID)
	if existing == nil {
		return 0
	}

	if val, ok := existing.(int64); ok {
		return val
	}

	return 0
}

func WithAccountID(ctx context.Context, accountID int64) context.Context {
	return context.WithValue(ctx, dtos.ContextKeyAccountID, accountID)
}
