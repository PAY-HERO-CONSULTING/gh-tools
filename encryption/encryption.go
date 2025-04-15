package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

var encryptionConfig config

type config struct {
	EncryptionKey []byte
	HMACKey       []byte
}

// Initialize encryption config with AES and HMAC keys
func NewConfig(secretKey, hmacKey string) {
	aesKey, err := convertKeyStringToBytes(secretKey)
	if err != nil {
		panic(err)
	}

	hmacKeyBytes, err := convertKeyStringToBytes(hmacKey)
	if err != nil {
		panic(err)
	}

	encryptionConfig = config{
		EncryptionKey: aesKey,
		HMACKey:       hmacKeyBytes,
	}
}

// Convert a hex string to bytes
func convertKeyStringToBytes(keyStr string) ([]byte, error) {
	return hex.DecodeString(keyStr)
}

// Encryption interface
type Encryption interface {
	Encrypt(ctx context.Context, plainText string) (string, error)
	Decrypt(ctx context.Context, encrypted string) (string, error)
}

type aesEncryption struct{}

// Factory function for AES encryption
func NewAESEncryption() Encryption {
	return &aesEncryption{}
}

// Encrypt using AES-GCM and HMAC-SHA256
func (e *aesEncryption) Encrypt(ctx context.Context, plainText string) (string, error) {
	block, err := aes.NewCipher(encryptionConfig.EncryptionKey)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt data using AES-GCM
	ciphertext := aead.Seal(nil, nonce, []byte(plainText), nil)

	// Generate HMAC-SHA256 for integrity
	h := hmac.New(sha256.New, encryptionConfig.HMACKey)
	h.Write(nonce)
	h.Write(ciphertext)
	hmacValue := h.Sum(nil)

	// Concatenate nonce + ciphertext + HMAC
	finalCiphertext := append(nonce, ciphertext...)
	finalCiphertext = append(finalCiphertext, hmacValue...)

	return hex.EncodeToString(finalCiphertext), nil
}

// Decrypt using AES-GCM and verify integrity with HMAC-SHA256
func (e *aesEncryption) Decrypt(ctx context.Context, encrypted string) (string, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", errors.New("failed to decode encrypted hex string")
	}

	block, err := aes.NewCipher(encryptionConfig.EncryptionKey)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	hmacSize := sha256.Size // HMAC-SHA256 produces a 32-byte digest

	if len(data) < nonceSize+hmacSize {
		return "", errors.New("invalid encrypted data length")
	}

	// Extract nonce, ciphertext, and HMAC
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize : len(data)-hmacSize]
	receivedHMAC := data[len(data)-hmacSize:]

	// Verify HMAC for integrity
	h := hmac.New(sha256.New, encryptionConfig.HMACKey)
	h.Write(nonce)
	h.Write(ciphertext)
	expectedHMAC := h.Sum(nil)

	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return "", errors.New("data integrity check failed: possible tampering detected")
	}

	// Decrypt using AES-GCM
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
