package moneydance

type Account struct {
	Id         string `json:"id"`
	ParentId   string `json:"parentid"`
	Name       string `json:"name"`
	Comment    string `json:"comment"`
	CurrencyId string `json:"currid"`
}
