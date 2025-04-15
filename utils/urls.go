package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func BuildPasswordResetURLV1(baseURL, resetKey string) (string, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/verification"

	query := parsedURL.Query()
	query.Set("reset_key", resetKey)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func BuildPasswordResetURL(baseURL, resetKey string) (string, error) {

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/" + resetKey

	return parsedURL.String(), nil
}

func BuildEmailVerificationURL(
	baseURL,
	userID,
	verificationCode string,
) (string, error) {

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure the path is properly formatted
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/" + url.PathEscape(userID) + "/" + url.PathEscape(verificationCode)

	return parsedURL.String(), nil
}
