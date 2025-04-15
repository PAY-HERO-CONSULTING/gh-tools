package utils

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateSalt() string {
	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 12)
}

func generateRandomString(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}
func randomStringFromSet(values ...string) string {
	n := len(values)
	if n == 0 {
		return ""
	}

	return values[rand.Intn(len(values))]
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomId() string {
	return uuid.New().String()
}

func GenerateRandomString(length int) string {
	return generateRandomString(lowerCaseLetters, length)
}

func GenerateRandomNumber(length int) string {
	return generateRandomString(digits, length)
}

func GenerateAvatarFileNameUUID() string {
	return generateRandomString(upperCaseLetters+lowerCaseLetters+digits, 8)
}

func GeneratePhoneActivationCode() string {
	return generateRandomString(digits, 5)
}

func GenerateAccountUUID() string {
	return generateRandomString(hexLetters+digits, 10)
}

func GenerateDocumentUUID() string {
	return generateRandomString(upperCaseLetters+lowerCaseLetters+digits, 8)
}
