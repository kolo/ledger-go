package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type balanceItem struct {
	account *account
	total   decimal.Decimal
}

func (bi *balanceItem) increase(amount decimal.Decimal) {
	bi.total = bi.total.Add(amount)
}

func (bi *balanceItem) decrease(amount decimal.Decimal) {
	bi.total = bi.total.Sub(amount)
}

type balance map[string]*balanceItem

func (b balance) update(r *record) {
	b.item(r.credit).decrease(r.amount)
	b.item(r.debit).increase(r.amount)
}

func (b balance) item(ac *account) *balanceItem {
	item, found := b[ac.name]
	if !found {
		item = &balanceItem{
			account: ac,
			total:   decimal.Zero,
		}
		b[ac.name] = item
	}

	return item
}

func balanceReport(rd recordReader) {
	balance := balance{}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		balance.update(r)
	}

	for _, bi := range balance {
		if bi.account.asset {
			fmt.Printf("%s: %v\n", bi.account.name, bi.total)
		}
	}
}
