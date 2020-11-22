package cmd

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type exportCmd struct {
	*baseCmd

	format string
}

func newExportCmd() *cobra.Command {
	export := &exportCmd{
		baseCmd: newBaseCmd(),
	}

	cmd := &cobra.Command{
		Use: "export",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return export.Run(cmd.OutOrStdout())
		},
	}

	export.addFlags(cmd.Flags())

	export.run = func(_ *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		switch export.format {
		case "csv":
			return exportToCSV(iter, stdout)
		case "sqlite":
			return exportToSQLite(iter, stdout)
		default:
			return errors.New("unknown export format")
		}
	}

	return cmd
}

func (c *exportCmd) addFlags(flags *pflag.FlagSet) {
	c.baseCmd.addFlags(flags)
	flags.StringVarP(&c.format, "format", "f", "csv", "set export format")
}

func exportToCSV(iter ledger.RecordIterator, output io.Writer) error {
	f, err := os.OpenFile("ledger.csv", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(output)
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

func exportToSQLite(iter ledger.RecordIterator, output io.Writer) error {
	return nil
}
