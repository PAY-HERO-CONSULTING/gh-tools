package jwt_manager

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Login struct {
	Username       string            `json:"username"`
	Email          string            `json:"email"`
	UserID         string            `json:"user_id"`
	IsAdmin        bool              `json:"is_admin"`
	SessionID      string            `json:"session_id"`
	Status         string            `json:"status"`
	Roles          map[string]string `json:"roles"`
	OrganizationID string            `json:"organization_id"`
}

func (l Login) ClaimsForAccessToken() AccessTokenClaims {

	return l.claimsForUser()
}

func (l Login) claimsForUser() AccessTokenClaims {
	return AccessTokenClaims{
		UserID:         l.UserID,
		Username:       l.Username,
		Roles:          l.Roles,
		Status:         l.Status,
		IsAdmin:        l.IsAdmin,
		SessionID:      l.SessionID,
		OrganizationID: l.OrganizationID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
