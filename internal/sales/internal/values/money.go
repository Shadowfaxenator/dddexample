package values

type Money struct {
	Currency  string
	Decimal   uint
	Precision uint
}

func NewMoney(currency string, Decimal uint, Precision uint) Money {
	return Money{Currency: currency, Decimal: Decimal, Precision: Precision}
}
