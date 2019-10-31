package cmd

import (
	"encoding/csv"
	"os"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	export := newBaseCmd()

	cmd := &cobra.Command{
		Use: "export",
		RunE: func(*cobra.Command, []string) error {
			return export.Run()
		},
	}

	export.addFlags(cmd.Flags())

	export.run = func(_ *ledger.Config, iter ledger.RecordIterator) error {
		return exportToCSV(iter)
	}

	return cmd
}

func exportToCSV(iter ledger.RecordIterator) error {
	f, err := os.OpenFile("ledger.csv", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for {
		r := iter.Next()
		if r == nil {
			break
		}

		if err = w.Write(r.ToArray()); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}
