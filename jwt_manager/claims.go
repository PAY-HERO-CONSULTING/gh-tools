package jwt_manager

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const ACCESS_TOKEN_DURATION = time.Hour * 1
const REFRESH_TOKEN_DURATION = time.Hour * 24 * 30 * 4

type RefreshTokenClaims struct {
	TokenType      string            `json:"token_type"`
	UserID         string            `json:"user_id"`
	AccountIDs     []string          `json:"account_ids,omitempty"`
	Username       string            `json:"username,omitempty"`
	IsAdmin        bool              `json:"is_admin"`
	Roles          map[string]string `json:"roles,omitempty"`
	Status         string            `json:"status"`
	SessionID      string            `json:"session_id"`
	IsPinSet       bool              `json:"is_pin_set,omitempty"`
	PhoneVerified  bool              `json:"phone_verified,omitempty"`
	EmailVefifed   bool              `json:"email_verified,omitempty"`
	Scope          string            `json:"scope,omitempty"`
	OrganizationID string            `json:"organization_id"`
	jwt.StandardClaims
}

type AccessTokenClaims struct {
	UserID         string            `json:"user_id"`
	AccountIDs     []string          `json:"account_ids,omitempty"`
	Username       string            `json:"username,omitempty"`
	IsAdmin        bool              `json:"is_admin"`
	Roles          map[string]string `json:"roles,omitempty"`
	Status         string            `json:"status"`
	SessionID      string            `json:"session_id"`
	IsPinSet       bool              `json:"is_pin_set,omitempty"`
	PhoneVerified  bool              `json:"phone_verified,omitempty"`
	EmailVefifed   bool              `json:"email_verified,omitempty"`
	Scope          string            `json:"scope,omitempty"`
	OrganizationID string            `json:"organization_id"`
	jwt.StandardClaims
}

func (c AccessTokenClaims) IsValidID(UserID string) bool {
	return c.UserID == UserID
}

func (c AccessTokenClaims) RefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType:      "refresh_token",
		UserID:         c.UserID,
		AccountIDs:     c.AccountIDs,
		Username:       c.Username,
		IsAdmin:        c.IsAdmin,
		Status:         c.Status,
		SessionID:      c.SessionID,
		IsPinSet:       c.IsPinSet,
		PhoneVerified:  c.PhoneVerified,
		EmailVefifed:   c.EmailVefifed,
		Scope:          c.Scope,
		Roles:          c.Roles,
		OrganizationID: c.OrganizationID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_DURATION).Unix(),
		},
	}
}

func (c RefreshTokenClaims) AccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		UserID:        c.UserID,
		AccountIDs:    c.AccountIDs,
		Username:      c.Username,
		Roles:         c.Roles,
		Status:        c.Status,
		SessionID:     c.SessionID,
		IsPinSet:      c.IsPinSet,
		PhoneVerified: c.PhoneVerified,
		EmailVefifed:  c.EmailVefifed,
		Scope:         c.Scope,
		IsAdmin:       c.IsAdmin,
		// RegisteredClaims: jwt.RegisteredClaims{
		// 	Issuer:    "payd.auth",
		// 	ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		// },
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}
