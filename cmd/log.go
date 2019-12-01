package cmd

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	log := newBaseCmd()

	cmd := &cobra.Command{
		Use: "log",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return log.Run(cmd.OutOrStdout())
		},
	}

	log.addFlags(cmd.Flags())

	log.run = func(_ *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		logRecords(iter, stdout)
		return nil
	}

	return cmd
}

func logRecords(iter ledger.RecordIterator, output io.Writer) {
	w := tabwriter.NewWriter(output, 0, 0, 2, ' ', 0)

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
