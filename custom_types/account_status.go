package custom_types

import (
	"database/sql/driver"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusPending  AccountStatus = "pending"
	AccountStatusBanned   AccountStatus = "banned"
)

func (s *AccountStatus) Scan(value interface{}) error {
	*s = AccountStatus(string(value.([]uint8)))
	return nil
}

func (s AccountStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s AccountStatus) String() string {
	return string(s)
}

func (s AccountStatus) IsActive() bool {
	return s == AccountStatusActive
}
