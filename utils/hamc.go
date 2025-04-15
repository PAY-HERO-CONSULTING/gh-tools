package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func VerifyHMAC(
	secret,
	timestamp,
	path,
	method,
	body,
	providedSignature string,
) bool {

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(timestamp))
	h.Write([]byte(path))
	h.Write([]byte(method))

	if body != "" {
		bodyHash := sha256.Sum256([]byte(body))
		bodyBase64 := base64.StdEncoding.EncodeToString(bodyHash[:])
		h.Write([]byte(bodyBase64))
	}

	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(providedSignature))
}
