package custom_types

import "database/sql/driver"

type ActivationChannel string

const (
	ActivationChannelEmail ActivationChannel = "email"
	ActivationChannelPhone ActivationChannel = "phone"
)

func (a *ActivationChannel) Scan(value interface{}) error {
	*a = ActivationChannel(value.(string))
	return nil
}

func (a ActivationChannel) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a ActivationChannel) String() string {
	return string(a)
}
