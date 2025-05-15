package dtos

import (
	"errors"
	"fmt"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
	"github.com/PAY-HERO-CONSULTING/gh-tools/serviceutils"
)

type Filter struct {
	AccountID           *string
	Category            string
	BusinessAccountID   *string
	ApplicationID       *string
	SerialNumber        *string
	AssignedAccountID   *string
	AssignedUserID      *string
	OrganizationID      *string
	CheckoutRequestID   *string
	LedgerType          *string
	FromTime            *time.Time
	ToTime              *time.Time
	Page                int
	Per                 int
	From                string
	ResetKey            string
	Source              *string
	To                  string
	Token               string
	Term                string
	UUID                string
	UserID              *string
	Deleted             *bool
	Status              string
	Type                string
	Valid               *bool
	Year                *string
	Username            *string
	Number              *string
	Email               *string
	WalletType          string
	TransactionCategory string
	SecurityQuestionId  string
	Provider            *string
	BillID              *string
	InvoiceID           *string
}

func (f *Filter) NoPagination() *Filter {
	return &Filter{
		From:                f.From,
		To:                  f.To,
		FromTime:            f.FromTime,
		ToTime:              f.ToTime,
		Term:                f.Term,
		UUID:                f.UUID,
		UserID:              f.UserID,
		Deleted:             f.Deleted,
		Status:              f.Status,
		Type:                f.Type,
		Token:               f.Token,
		Valid:               f.Valid,
		AccountID:           f.AccountID,
		BusinessAccountID:   f.BusinessAccountID,
		AssignedAccountID:   f.AssignedAccountID,
		AssignedUserID:      f.AssignedUserID,
		OrganizationID:      f.OrganizationID,
		LedgerType:          f.LedgerType,
		CheckoutRequestID:   f.CheckoutRequestID,
		Source:              f.Source,
		Category:            f.Category,
		Year:                f.Year,
		Username:            f.Username,
		Number:              f.Number,
		Email:               f.Email,
		WalletType:          f.WalletType,
		TransactionCategory: f.TransactionCategory,
		Provider:            f.Provider,
		BillID:              f.BillID,
		InvoiceID:           f.InvoiceID,
	}
}

func (f *Filter) ConvertTime() error {
	if f.From == "" || f.To == "" {
		return errors.New("time_filter: from or to filter time is empty")
	}

	fromTime, err := serviceutils.ParseTime(f.From)
	if err != nil {
		return fmt.Errorf("parse from time [%v], err [%v]", f.From, err)
	}

	f.FromTime = null.NullValue(fromTime)

	toTime, err := serviceutils.ParseTime(f.To)
	if err != nil {
		return fmt.Errorf("parse to time [%v], err [%v]", f.To, err)
	}

	f.ToTime = null.NullValue(toTime)

	return nil
}

func (f *Filter) TimeFilterSet() bool {
	return f.From != "" && f.To != ""
}

func (f *Filter) SetValid(val bool) {
	f.Valid = null.NullValue(val)
}
