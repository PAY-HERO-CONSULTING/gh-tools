package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const minPasswordLength = 8

var invalidPasswordRegex = regexp.MustCompile("[[:^graph:]]")

func GeneratePasswordHashV2(password string) string {
	salt := GenerateSalt()
	keyLength := 32
	iterations := 600000
	encodedPassword := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New)
	encoded := base64.StdEncoding.EncodeToString(encodedPassword)
	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iterations, salt, encoded)
}

func MatchPassword(passwordHash, password string) error {
	parts := strings.Split(passwordHash, "$")
	if parts == nil || len(parts) != 4 {
		return errors.New("invalid password hash")
	}

	if parts[0] != "pbkdf2_sha256" {
		return errors.New("invalid password hash")
	}

	iter, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	salt := []byte(parts[2])
	hash := parts[3]

	b, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return err
	}

	dk := pbkdf2.Key([]byte(password), salt, iter, sha256.Size, sha256.New)

	if !bytes.Equal(b, dk) {
		return errors.New("invalid password")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(strings.TrimSpace(password)) < minPasswordLength {
		return apperrs.NewAuthorization(
			fmt.Sprintf("Failed to validate password of length %d", len(password)),
		)
	}

	loc := invalidPasswordRegex.FindStringIndex(password)
	if loc != nil {
		return apperrs.NewAuthorization(
			"Failed to validate password with invalid characters",
		)
	}

	return nil
}

func CheckMatchingPasswords(password, passwordConfirmation string) error {
	if password != passwordConfirmation {
		return apperrs.NewAuthorization(
			"passwords do not match",
		)
	}

	return nil
}

func GenerateApiKeyPasswordHash(secret string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(secret), 6)
	return string(b), err
}

func CheckApiKeyPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func PasswordHashCost(passwordHash string) (int, error) {
	return bcrypt.Cost([]byte(passwordHash))
}
