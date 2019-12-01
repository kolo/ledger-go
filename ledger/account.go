package ledger

import "github.com/shopspring/decimal"

type Account struct {
	Name  string
	Asset bool
}

func NewAccount(name string, asset bool) *Account {
	return &Account{
		Name:  name,
		Asset: asset,
	}
}

func (account *Account) String() string {
	return account.Name
}

type accountBalance struct {
	account *Account
	amount  decimal.Decimal
}

func newAccountBalance(account *Account) *accountBalance {
	return &accountBalance{
		account: account,
		amount:  decimal.Zero,
	}
}

func (b *accountBalance) increase(amount decimal.Decimal) {
	b.amount = b.amount.Add(amount)
}

func (b *accountBalance) decrease(amount decimal.Decimal) {
	b.amount = b.amount.Sub(amount)
}
