package custom_types

import "database/sql/driver"

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusAccepted KYCStatus = "accepted"
	KYCStatusRejected KYCStatus = "rejected"
)

func (k *KYCStatus) Scan(value any) error {
	*k = KYCStatus(string(value.([]uint8)))
	return nil
}

func (k KYCStatus) Value() (driver.Value, error) {
	return k.String(), nil
}

func (k KYCStatus) String() string {
	return string(k)
}
