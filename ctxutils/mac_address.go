package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
)

func MacAddress(ctx context.Context) string {
	existing := ctx.Value(dtos.ContextKeyMacAddress)
	if existing == nil {
		return ""
	}

	return existing.(string)
}

func WithMacAddressAddress(ctx context.Context, macAddress string) context.Context {
	return context.WithValue(ctx, dtos.ContextKeyMacAddress, macAddress)
}
