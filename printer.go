package main

import (
	"fmt"
	"os"
	"sort"
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
