package ledger

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/shopspring/decimal"
)

func ExpensesReport(iter RecordIterator) {
	iter = NewFilteredIterator(iter, filterExpenses)

	credits, debits := map[string]struct{}{}, map[string]struct{}{}
	expenses := map[[2]string]decimal.Decimal{}

	update := func(r *Record) {
		key := [2]string{r.Credit.Name, r.Debit.Name}
		if _, found := expenses[key]; !found {
			expenses[key] = decimal.Zero
		}

		if _, found := credits[r.Credit.Name]; !found {
			credits[r.Credit.Name] = struct{}{}
		}

		if _, found := debits[r.Debit.Name]; !found {
			debits[r.Debit.Name] = struct{}{}
		}

		expenses[key] = expenses[key].Add(r.Amount)
	}

	for {
		r := iter.Next()
		if r == nil {
			break
		}

		update(r)
	}

	// printing
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	columns := []string{}
	totals := map[string]decimal.Decimal{}
	for credit := range credits {
		columns = append(columns, credit)
		totals[credit] = decimal.Zero
	}
	sort.Strings(columns)

	rows := []string{}
	for debit := range debits {
		rows = append(rows, debit)
	}
	sort.Strings(rows)

	// header
	fmt.Fprintf(w, "\t%v\t=\n", strings.Join(columns, "\t"))

	// body
	for _, debit := range rows {
		v := []string{debit}
		t := decimal.Zero

		for _, credit := range columns {
			key := [2]string{credit, debit}
			spend := expenses[key]
			if spend.Equals(decimal.Zero) {
				v = append(v, "-")
			} else {
				v = append(v, spend.StringFixed(2))
				t = t.Add(spend)
				totals[credit] = totals[credit].Add(spend)
			}
		}

		fmt.Fprintf(w, "%s\t%s\n", strings.Join(v, "\t"), t.StringFixed(2))
	}

	// footer
	fmt.Fprintf(w, "=")
	total := decimal.Zero
	for _, column := range columns {
		total = total.Add(totals[column])
		fmt.Fprintf(w, "\t%s", totals[column].StringFixed(2))
	}
	fmt.Fprintf(w, "\t%s", total.StringFixed(2))
	fmt.Fprintf(w, "\n")

	w.Flush()
}

func filterExpenses(r *Record) *Record {
	if r.RecordType() != recordTypeExpense {
		return nil
	}

	return r
}
