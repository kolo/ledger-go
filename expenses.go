package main

import (
	"fmt"
	"time"

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

	t, _ := time.Parse(time.RFC3339, "2018-05-21T00:00:00+02:00")
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
		fmt.Printf("%s: %v\n", bi.account.name, bi.total.StringFixed(2))
	}
}
