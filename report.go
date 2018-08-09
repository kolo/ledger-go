package main

import "github.com/shopspring/decimal"

type reportItem struct {
	account *account
	total   decimal.Decimal
}

func (ri *reportItem) increase(amount decimal.Decimal) {
	ri.total = ri.total.Add(amount)
}

func (ri *reportItem) decrease(amount decimal.Decimal) {
	ri.total = ri.total.Sub(amount)
}

type report map[string]*reportItem

func (r report) update(rec *record) {
	r.item(rec.credit).decrease(rec.amount)
	r.item(rec.debit).increase(rec.amount)
}

func (r report) item(ac *account) *reportItem {
	item, found := r[ac.name]
	if !found {
		item = &reportItem{
			account: ac,
			total:   decimal.Zero,
		}
		r[ac.name] = item
	}

	return item
}

func (r report) total() decimal.Decimal {
	total := decimal.Zero
	for _, ri := range r {
		total = total.Add(ri.total)
	}

	return total
}
