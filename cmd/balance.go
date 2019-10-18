package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var unixEpoch = time.Unix(0, 0)

type balanceCmd struct {
	ledgerDir string
	period    *dateRangeFlags
}

func newBalanceCmd() *cobra.Command {
	balance := &balanceCmd{
		period: newDateRangeFlags(unixEpoch, time.Now()),
	}

	cmd := &cobra.Command{
		Use: "balance",
		RunE: func(*cobra.Command, []string) error {
			return balance.Execute()
		},
	}

	balance.addFlags(cmd.Flags())

	return cmd
}

func (c *balanceCmd) addFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&c.ledgerDir, "ledger-dir", "", os.Getenv("LEDGER_DIR"), "set ledger directory")
	c.period.addFlags(flags)
}

func (c *balanceCmd) Execute() error {
	cfg, err := ledger.LoadConfig(filepath.Join(c.ledgerDir, "config.json"))
	if err != nil {
		return err
	}

	ledger.PrintBalanceReport(
		ledger.NewLedgerIterator(cfg.Assets, c.ledgerDir),
		cfg.Assets,
	)

	return nil
}
