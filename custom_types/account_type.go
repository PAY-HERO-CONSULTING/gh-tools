package custom_types

import "database/sql/driver"

type AccountsType string

const (
	FixedAssetAccountsType        AccountsType = "fixed_asset"
	CurrentAssetAccountsType      AccountsType = "current_asset"
	LiabilityAccountsType         AccountsType = "liability"
	CurrentLiabilityAccountsType  AccountsType = "current_liability"
	ExpenseAccountsType           AccountsType = "expense"
	CostOfGoodsAccountsType       AccountsType = "cost_of_goods"
	OtherExpenseAccountsType      AccountsType = "other_expense"
	RevenueAccountsType           AccountsType = "revenue"
	AccountReceivableAccountsType AccountsType = "accounts_receivable"
)

func (a *AccountsType) Scan(value interface{}) error {
	*a = AccountsType(string(value.([]uint8)))
	return nil
}

func (s AccountsType) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s AccountsType) String() string {
	return string(s)
}
