package appvalidator

func init() {
	Default().RegisterAlias("currencyIsoCode", "len=3,alpha,uppercase")
}
