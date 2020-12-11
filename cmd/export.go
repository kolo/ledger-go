package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

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

	export.run = func(config *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		switch export.format {
		case "csv":
			return exportToCSV(iter, stdout)
		case "ldg":
			return exportToLdg(config, iter, stdout)
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

func exportToLdg(config *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
	f, err := os.OpenFile("journal.ldg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	var lastDate time.Time

	flushPairing := func(p []*ledger.Record) {
		if len(p) == 0 {
			return
		}
		if len(p) == 1 {
			f.WriteString(fmt.Sprintf("  %s,%s,%s\n", p[0].Credit.Name, p[0].Debit.Name, p[0].Amount.StringFixed(2)))
			return
		}
		f.WriteString(fmt.Sprintf("  %s {\n", p[0].Credit.Name))
		for _, pp := range p {
			f.WriteString(fmt.Sprintf("    %s,%s\n", pp.Debit.Name, pp.Amount.StringFixed(2)))
		}
		f.WriteString("  }\n")
	}

	flushRecord := func(r []*ledger.Record) {
		if len(r) == 0 {
			return
		}
		f.WriteString(fmt.Sprintf("%s {\n", r[0].Date.Format("2006-01-02")))
		var i int
		var lastCredit string
		for j, rr := range r {
			if rr.Credit.Name != lastCredit {
				flushPairing(r[i:j])
				lastCredit = rr.Credit.Name
				i = j
			}
		}
		flushPairing(r[i:])
		f.WriteString("}\n")
	}

	record := []*ledger.Record{}
	for {
		r := iter.Next()
		if r == nil {
			break
		}

		if !r.Date.Equal(lastDate) {
			lastDate = r.Date
			flushRecord(record)
			record = record[:0]
			record = append(record, r)
		} else {
			record = append(record, r)
		}
	}
	return nil
}
