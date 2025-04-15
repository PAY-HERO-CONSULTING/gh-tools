package custom_types

import (
	"strings"

	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
)

type PhoneNumber struct {
	DialCode *string `json:"dial_code"`
	Number   *string `json:"number"`
}

func (p PhoneNumber) Phone() string {
	return strings.TrimSpace(null.ValueFromNull(p.DialCode) + null.ValueFromNull(p.Number))
}
