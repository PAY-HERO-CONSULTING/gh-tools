package utils

func GeneratePasswordResetKey() string {
	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 16)
}
