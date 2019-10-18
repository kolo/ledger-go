package ledger

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type balanceReport struct {
	balances map[string]*accountBalance
}

func newBalanceReport(assets []*Account) *balanceReport {
	report := &balanceReport{
		balances: map[string]*accountBalance{},
	}

	for _, asset := range assets {
		report.balances[asset.Name] = newAccountBalance(asset)
	}

	return report
}

func (report *balanceReport) update(r *Record) {
	if balance, ok := report.balances[r.Credit.Name]; ok {
		balance.decrease(r.Amount)
	}
	if balance, ok := report.balances[r.Debit.Name]; ok {
		balance.increase(r.Amount)
	}
}

func (report *balanceReport) total() decimal.Decimal {
	total := decimal.Zero
	for _, balance := range report.balances {
		total = total.Add(balance.amount)
	}

	return total
}

func PrintBalanceReport(iter RecordIterator, assets []*Account) {
	report := newBalanceReport(assets)

	for {
		r := iter.Next()
		if r == nil {
			break
		}

		report.update(r)
	}

	fmt.Println(report.total())
}
