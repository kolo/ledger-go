package main

import (
	"fmt"
	"time"

	"github.com/kolo/ledger-go/datetime"
	"github.com/shopspring/decimal"
)

type expenses map[string]*balanceItem

func (e expenses) update(r *record) {
	if r.recordType() != recordTypeExpense {
		return
	}

	if bi, found := e[r.credit.name]; found {
		bi.increase(r.amount)
	}
}

func (e expenses) total() decimal.Decimal {
	total := decimal.Zero
	for _, bi := range e {
		total = total.Add(bi.total)
	}

	return total
}

func expensesReport(rd recordReader, assets []string) {
	expenses := expenses{}
	for _, asset := range assets {
		expenses[asset] = &balanceItem{
			account: &account{
				name:  asset,
				asset: true,
			},
			total: decimal.Zero,
		}
	}

	t := datetime.BeginningOfWeek(time.Now())
	for {
		r := rd.Next()
		if r == nil {
			break
		}

		if r.recordedAt.After(t) {
			expenses.update(r)
		}
	}

	for _, bi := range expenses {
		fmt.Printf("%5s: %6s\n", bi.account.name, bi.total.StringFixed(2))
	}
	fmt.Printf("Total: %v\n", expenses.total())
}
