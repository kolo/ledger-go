package ledger

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func LogReport(iter RecordIterator) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for {
		r := iter.Next()
		if r == nil {
			break
		}

		date := r.Date.Format(ISO8601Date)
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			date,
			r.Credit,
			r.Debit,
			r.formatAmount(),
		)

	}

	w.Flush()
}
