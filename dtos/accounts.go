package dtos

type UpdateAccountRequest struct {
	AccountID     string  `json:"account_id"`
	AccountNumber *string `json:"account_number"`
}
