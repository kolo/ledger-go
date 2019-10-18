package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/shopspring/decimal"
)

func printWeeklyReport(r *weeklyReport) {
	// Sort assets
	assets := make([]string, len(r.assets))
	copy(assets, r.assets)
	sort.Strings(assets)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	fmt.Fprintf(w, "\t%s\tTotal\n", strings.Join(assets, "\t"))

	// Print report body
	totals := make([]decimal.Decimal, len(assets))
	for ri := r.head; ri != nil; ri = ri.next {
		subtotals := []string{}
		total := decimal.Zero
		for i, asset := range assets {
			total = total.Add(ri.report[asset].total)
			totals[i] = totals[i].Add(ri.report[asset].total)
			subtotals = append(subtotals, ri.report[asset].total.StringFixed(2))
		}
		fmt.Fprintf(w, "%v\t%s\t%s\n", ri.week, strings.Join(subtotals, "\t"), total.StringFixed(2))
	}

	// Print footer
	sum := decimal.Zero
	stringTotals := make([]string, len(totals))
	for i, total := range totals {
		sum = sum.Add(total)
		stringTotals[i] = total.StringFixed(2)
	}
	fmt.Fprintf(w, "-----\t%s\n", strings.Join(make([]string, len(totals)), "\t"))
	fmt.Fprintf(w, "Total\t%s\t%s\n", strings.Join(stringTotals, "\t"), sum.StringFixed(2))

	w.Flush()
}
