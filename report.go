package main

import "github.com/shopspring/decimal"

type reportItem struct {
	account *account
	total   decimal.Decimal
}

func (bi *reportItem) increase(amount decimal.Decimal) {
	bi.total = bi.total.Add(amount)
}

func (bi *reportItem) decrease(amount decimal.Decimal) {
	bi.total = bi.total.Sub(amount)
}

type report map[string]*reportItem

func (b report) update(r *record) {
	b.item(r.credit).decrease(r.amount)
	b.item(r.debit).increase(r.amount)
}

func (b report) item(ac *account) *reportItem {
	item, found := b[ac.name]
	if !found {
		item = &reportItem{
			account: ac,
			total:   decimal.Zero,
		}
		b[ac.name] = item
	}

	return item
}
