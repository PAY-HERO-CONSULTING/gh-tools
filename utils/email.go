package utils

import (
	"errors"
	"regexp"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
)

var (
	emailRe    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,10}$`)
	htmlTagsRe = regexp.MustCompile(`<(.|\n)*?>`)
)

func ValidateEmail(email string) error {
	if !emailRe.MatchString(email) {
		return apperrs.Wrap(
			errors.New("invalid email format"),
		)
	}
	return nil
}

func StripHtmlTags(email string) string {
	return htmlTagsRe.ReplaceAllString(email, "")
}
