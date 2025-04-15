package custom_types

import (
	"database/sql/driver"
)

type UserStatus string

const (
	UserStatusVerified UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
	UserStatusBanned   UserStatus = "banned"
)

// Scan implements the Scanner interface.
func (s *UserStatus) Scan(value interface{}) error {
	*s = UserStatus(string(value.([]uint8)))
	return nil
}

// Value implements the driver Valuer interface.
func (s UserStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsActive() bool {
	return s == UserStatusVerified
}

func (s UserStatus) IsUnverified() bool {
	return s == UserStatusPending
}
