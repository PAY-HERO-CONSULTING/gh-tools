package dtos

type ContextKey string

const (
	ContextKeyAccountID         ContextKey = "accountID"
	ContextKeyIpAddress         ContextKey = "ipAddress"
	ContextKeyLang              ContextKey = "lang"
	ContextKeyLeadingWhitespace ContextKey = "leadingWhitespace"
	ContextKeyMacAddress        ContextKey = "macAddress"
	ContextKeyRepeatRecipients  ContextKey = "repeatRecipientsKey"
	ContextKeyRequestID         ContextKey = "requestId"
	ContextKeyTokenInfo         ContextKey = "tokenInfo"
	ContextKeyUserAgent         ContextKey = "userAgent"
)
