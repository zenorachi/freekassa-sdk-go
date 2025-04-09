package freekassa

type Currency string

const (
	RUB Currency = "RUB"
	USD Currency = "USD"
	EUR Currency = "EUR"
	UAH Currency = "UAH"
	KZT Currency = "KZT"
)

func (c Currency) String() string {
	return string(c)
}
