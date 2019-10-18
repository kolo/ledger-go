package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	log := newBaseCmd()

	cmd := &cobra.Command{
		Use: "log",
		RunE: func(*cobra.Command, []string) error {
			return log.Run()
		},
	}

	log.addFlags(cmd.Flags())

	log.run = func(_ *ledger.Config, iter ledger.RecordIterator) error {
		logRecords(iter)
		return nil
	}

	return cmd
}

func logRecords(iter ledger.RecordIterator) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for {
		r := iter.Next()
		if r == nil {
			break
		}

		date := r.Date.Format(ledger.ISO8601Date)
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			date,
			r.Credit,
			r.Debit,
			r.FormatAmount(),
		)

	}

	w.Flush()
}
