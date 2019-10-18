package ledger

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

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

func (report *balanceReport) print() {
	var accounts []string
	for name := range report.balances {
		accounts = append(accounts, name)
	}
	sort.Strings(accounts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, account := range accounts {
		balance := report.balances[account]
		if balance.amount.Equal(decimal.Zero) {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\n", balance.account.Name, balance.amount.StringFixed(2))
	}
	fmt.Fprintln(w, "-----\t")
	fmt.Fprintf(w, "Total\t%s\n", report.total().StringFixed(2))

	w.Flush()
}

func BalanceReport(iter RecordIterator, assets []*Account) {
	report := newBalanceReport(assets)
	for {
		r := iter.Next()
		if r == nil {
			break
		}

		report.update(r)
	}

	report.print()
}
