package custom_types

import "database/sql/driver"

type OrganizationStatus string

const (
	OrganizationStatusActive   OrganizationStatus = "active"
	OrganizationStatusInactive OrganizationStatus = "inactive"
	OrganizationStatusBanned   OrganizationStatus = "banned"
)

func (s *OrganizationStatus) Scan(value any) error {
	*s = OrganizationStatus(string(value.([]uint8)))
	return nil
}

func (s OrganizationStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s OrganizationStatus) String() string {
	return string(s)
}
