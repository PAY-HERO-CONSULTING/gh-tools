package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Contains(value, subValue string) bool {
	return strings.Contains(value, subValue)
}

func Title(str string) string {
	return cases.Title(language.Und).String(str)
}
