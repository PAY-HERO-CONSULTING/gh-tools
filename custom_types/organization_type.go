package custom_types

import (
	"database/sql/driver"
)

type OrganizationType string

const (
	OrganizationTypeIndividual OrganizationType = "individual"
	OrganizationTypeBusiness   OrganizationType = "business"
)

// Scan implements the Scanner interface
func (s *OrganizationType) Scan(value interface{}) error {
	*s = OrganizationType(string(value.([]uint8)))
	return nil
}

// Value implement the driver Valuer interface.
func (s OrganizationType) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s OrganizationType) String() string {
	return string(s)
}
