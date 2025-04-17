package jwt_manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/models"
	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	siginingKey []byte
	keyFunc     func(*jwt.Token) (interface{}, error)
}

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken(secret string) (string, error) {
	signedString, err := t.token.SignedString([]byte(secret))
	if err != nil {
		return "", apperrs.NewUnexpectedError(err.Error()).AddLogMessage(
			"cannot generate access token",
		)
	}

	return signedString, nil
}

func (t AuthToken) NewRefreshToken(
	secret string,
) (string, error) {
	c := t.token.Claims.(AccessTokenClaims)

	refreshClaims := c.RefreshTokenClaims()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", apperrs.NewUnexpectedError(err.Error()).AddLogMessage("Failed while signing refresh token")
	}

	return signedString, nil
}

func NewAuthToken(
	claims AccessTokenClaims,
) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return AuthToken{
		token: token,
	}
}

func NewAccessTokenFromRefreshToken(
	hmacSecretKey string,
	refreshToken string,
) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSecretKey), nil
	})
	if err != nil {
		return "", apperrs.NewAuthorization("invalid or expired refresh token")
	}

	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)

	return authToken.NewAccessToken(hmacSecretKey)
}

func NewJWTManager(
	signingKey string,
) *JWTManager {
	jwtHandler := &JWTManager{
		siginingKey: []byte(signingKey),
	}

	jwtHandler.keyFunc = func(*jwt.Token) (interface{}, error) {
		return jwtHandler.siginingKey, nil
	}

	return jwtHandler
}

func (m *JWTManager) Generate(
	hmacSecretKey string,
	organizationID string,
	sessionID string,
	user *models.User,
) (string, error) {
	session := &Session{
		Username:       user.UserName,
		IsAdmin:        user.IsAdmin,
		UserID:         user.ID,
		Status:         user.Status.String(),
		IsPinSet:       user.IsPinSet,
		PhoneVerified:  user.PhoneVerified,
		EmailVefifed:   user.EmailVefifed,
		Scope:          "gh-system",
		SessionID:      sessionID,
		OrganizationID: organizationID,
	}

	claims := session.ClaimsForAccessToken()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(hmacSecretKey))
}

func (m *JWTManager) Verify(
	hmacSecretKey string,
	accessToken string,
) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(hmacSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token err = [%+v]", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims error")
	}

	return claims, nil
}

func (m *JWTManager) TokenInfo(tokenValue string) (*models.TokenInfo, error) {
	token, err := jwt.Parse(tokenValue, m.keyFunc)
	if err != nil {
		return &models.TokenInfo{}, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &models.TokenInfo{}, nil
	}

	exp := m.getInt64(mapClaims, "exp")
	isAdmin := m.getBool(mapClaims, "is_admin")
	status := m.getString(mapClaims, "status")
	userID := m.getString(mapClaims, "user_id")
	accountIDs := m.getStringSlice(mapClaims, "account_ids")
	roles := m.getStringMap(mapClaims, "roles")
	isPinSet := m.getBool(mapClaims, "is_pin_set")
	phoneVerified := m.getBool(mapClaims, "phone_verified")
	emailVerified := m.getBool(mapClaims, "email_verified")
	scope := m.getString(mapClaims, "scope")
	sessionID := m.getString(mapClaims, "session_id")
	username := m.getString(mapClaims, "username")

	tokenInfo := &models.TokenInfo{
		AccountIDs:    accountIDs,
		Exp:           time.Unix(exp, 0),
		IsAdmin:       isAdmin,
		Status:        status,
		UserID:        userID,
		Roles:         roles,
		Scope:         scope,
		IsPinSet:      isPinSet,
		EmailVefifed:  emailVerified,
		PhoneVerified: phoneVerified,
		SessionID:     sessionID,
		Username:      username,
	}

	return tokenInfo, nil
}

func (m *JWTManager) getInt64(mapClaims jwt.MapClaims, key string) int64 {
	switch val := mapClaims[key].(type) {
	case float64:
		return int64(val)
	case json.Number:
		v, _ := val.Int64()
		return v
	}
	return 0
}

func (m *JWTManager) getString(mapClaims jwt.MapClaims, key string) string {
	switch val := mapClaims[key].(type) {
	case interface{}:
		return val.(string)
	}
	return ""
}

func (h *JWTManager) getBool(mapClaims jwt.MapClaims, key string) bool {
	switch val := mapClaims[key].(type) {
	case interface{}:
		return val.(bool)
	}

	return false
}

func (h *JWTManager) getStringSlice(mapClaims jwt.MapClaims, key string) []string {
	stringSlices := make([]string, 0)

	switch val := mapClaims[key].(type) {
	case []interface{}:
		for _, v := range val {
			stringSlices = append(stringSlices, v.(string))
		}
		return stringSlices
	}

	return []string{}
}

func (h *JWTManager) getStringMap(mapClaims jwt.MapClaims, key string) map[string]string {
	int64StringMap := make(map[string]string, 0)

	switch val := mapClaims[key].(type) {
	case map[string]interface{}:
		for id, role := range val {
			int64StringMap[id] = role.(string)
		}
		return int64StringMap
	}

	return map[string]string{}
}
