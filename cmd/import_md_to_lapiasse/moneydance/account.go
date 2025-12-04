package moneydance

type Account struct {
	Id                   string      `json:"id"`
	AcctId               string      `json:"acctid"`
	OldId                string      `json:"old_id"`
	ParentId             string      `json:"parentid"`
	Name                 string      `json:"name"`
	Type                 AccountType `json:"type"`
	Comment              string      `json:"comment"`
	CurrencyId           string      `json:"currid"`
	IsInactive           Boolean     `json:"is_inactive"`
	CreationDate         UnixDate    `json:"creation_date"`
	InitialBalance       int64       `json:"sbal,string"`
	DefaultCategoryOldId string      `json:"default_category"`
	ExpiryYear           int         `json:"exp_year,string"`
	ExpiryMonth          int         `json:"exp_month,string"`
}

type AccountType string

const (
	AccountTypeAsset      AccountType = "a"
	AccountTypeBank       AccountType = "b"
	AccountTypeCard       AccountType = "c"
	AccountTypeExpense    AccountType = "e"
	AccountTypeIncome     AccountType = "i"
	AccountTypeLiability  AccountType = "l"
	AccountTypeLoan       AccountType = "o"
	AccountTypeInvestment AccountType = "v"
)
