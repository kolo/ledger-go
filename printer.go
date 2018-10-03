package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/shopspring/decimal"
)

func printReport(r report) {
	// Sort accounts
	accounts := []string{}
	for k := range r {
		accounts = append(accounts, k)
	}
	sort.Strings(accounts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, account := range accounts {
		ri := r[account]
		if !ri.total.Equal(decimal.Zero) {
			fmt.Fprintf(w, "%s\t%s\n", ri.account.name, ri.total.StringFixed(2))
		}
	}
	fmt.Fprintln(w, "-----\t")
	fmt.Fprintf(w, "Total\t%s\n", r.total().StringFixed(2))

	w.Flush()
}

func printWeeklyReport(r *weeklyReport) {
	// Sort assets
	assets := make([]string, len(r.assets))
	copy(assets, r.assets)
	sort.Strings(assets)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintf(w, "\t%s\n", strings.Join(assets, "\t"))

	// Print report body
	totals := make([]decimal.Decimal, len(assets))
	for ri := r.head; ri != nil; ri = ri.next {
		subtotals := []string{}
		for i, asset := range assets {
			totals[i] = totals[i].Add(ri.report[asset].total)
			subtotals = append(subtotals, ri.report[asset].total.StringFixed(2))
		}
		fmt.Fprintf(w, "%v\t%s\n", ri.id, strings.Join(subtotals, "\t"))
	}

	// Print footer
	stringTotals := make([]string, len(totals))
	for i, total := range totals {
		stringTotals[i] = total.StringFixed(2)
	}
	fmt.Fprintf(w, "-----\t%s\n", strings.Join(make([]string, len(totals)), "\t"))
	fmt.Fprintf(w, "Total\t%s\n", strings.Join(stringTotals, "\t"))

	w.Flush()
}
