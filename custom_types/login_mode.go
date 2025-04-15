package custom_types

import "database/sql/driver"

type LoginMode string

const (
	LoginModeEmail    LoginMode = "email"
	LoginModePhone    LoginMode = "phone"
	LoginModeUsername LoginMode = "username"
)

func (l *LoginMode) Scan(value interface{}) error {
	*l = LoginMode(string(value.([]uint8)))
	return nil
}

func (l LoginMode) Value() (driver.Value, error) {
	return l.String(), nil
}

func (l LoginMode) String() string {
	return string(l)
}
