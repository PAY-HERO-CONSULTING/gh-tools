package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() (string, error) {
	codeLength := 4

	return GenerateRandomInt(codeLength)
}

func GenerateRandomInt(codeLength int) (string, error) {
	code := ""

	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code += string(digits[num.Int64()])
	}

	return code, nil
}
