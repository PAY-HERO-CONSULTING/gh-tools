package jwt_manager

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Session struct {
	Username       string `json:"username"`
	UserID         string `json:"user_id"`
	IsAdmin        bool   `json:"is_admin"`
	SessionID      string `json:"session_id"`
	Status         string `json:"status"`
	IsPinSet       bool   `json:"is_pin_set"`
	PhoneVerified  bool   `json:"phone_verified"`
	EmailVefifed   bool   `json:"email_verified"`
	OrganizationID string `json:"organization_id"`
	Scope          string `json:"scope"`
}

func (s Session) ClaimsForAccessToken() AccessTokenClaims {

	return s.claimsForUser()
}

func (s Session) claimsForUser() AccessTokenClaims {
	return AccessTokenClaims{
		UserID:         s.UserID,
		Username:       s.Username,
		Status:         s.Status,
		IsAdmin:        s.IsAdmin,
		SessionID:      s.SessionID,
		Scope:          s.Scope,
		IsPinSet:       s.IsPinSet,
		PhoneVerified:  s.PhoneVerified,
		EmailVefifed:   s.EmailVefifed,
		OrganizationID: s.OrganizationID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
