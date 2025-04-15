package ctxutils

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/custom_types"
	"github.com/PAY-HERO-CONSULTING/gh-tools/models"
)

func WithTokenInfo(ctx context.Context, tokenInfo *models.TokenInfo) context.Context {
	return context.WithValue(ctx, custom_types.ContextKeyTokenInfo, tokenInfo)
}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, custom_types.ContextKeyToken, token)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	tokenInfo := TokenInfo(ctx)
	tokenInfo.UserID = userID
	return WithTokenInfo(ctx, tokenInfo)
}

func WithWorker(ctx context.Context) context.Context {
	tokenInfo := TokenInfo(ctx)
	tokenInfo.IsWorker = true
	return WithTokenInfo(ctx, tokenInfo)
}

func TokenInfo(ctx context.Context) *models.TokenInfo {
	existing := ctx.Value(custom_types.ContextKeyTokenInfo)
	if existing == nil {
		return &models.TokenInfo{}
	}

	tokenInfo, ok := existing.(*models.TokenInfo)
	if !ok {
		return &models.TokenInfo{}
	}

	return tokenInfo
}

func UserID(ctx context.Context) string {
	return TokenInfo(ctx).UserID
}

func DefaultAccount(ctx context.Context) string {
	return TokenInfo(ctx).AccountIDs[0]
}

func AccountIDs(ctx context.Context) []string {
	return TokenInfo(ctx).AccountIDs
}

func AccountUserRoles(ctx context.Context) map[string]string {
	return TokenInfo(ctx).Roles
}

func AccessToken(ctx context.Context) string {
	return TokenInfo(ctx).Token
}
