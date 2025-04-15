package custom_types

type ContextKey string

const (
	ContextKeyAccountIDs ContextKey = "accountIDs"
	ContextKeyIpAddress  ContextKey = "ipAddress"
	ContextKeyRequestID  ContextKey = "requestId"
	ContextKeyTokenInfo  ContextKey = "tokenInfo"
	ContextKeyToken      ContextKey = "token"
	ContextKeyUserAgent  ContextKey = "userAgent"
)
