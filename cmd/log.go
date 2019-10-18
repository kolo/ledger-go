package cmd

import (
	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

type logCmd struct {
	*baseCmd
}

func newLogCmd() *cobra.Command {
	log := &logCmd{
		baseCmd: newBaseCmd(),
	}

	cmd := &cobra.Command{
		Use: "log",
		RunE: func(*cobra.Command, []string) error {
			return log.Run()
		},
	}

	log.addFlags(cmd.Flags())

	log.run = func(_ *ledger.Config, iter ledger.RecordIterator) error {
		ledger.LogReport(iter)
		return nil
	}

	return cmd
}
