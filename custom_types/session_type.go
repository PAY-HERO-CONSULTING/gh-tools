package custom_types

import "database/sql/driver"

type SessionType string

const (
	SessionTypeJWT    SessionType = "jwt"
	SessionTypeCookie SessionType = "cookie"
)

func (s *SessionType) Scan(value any) error {
	*s = SessionType(string(value.([]uint8)))
	return nil
}

func (s SessionType) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s SessionType) String() string {
	return string(s)
}
