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
	accounts *accountFlags
	period   *dateRangeFlags

	ledgerDir string
}

func newBalanceCmd() *cobra.Command {
	balance := &balanceCmd{
		accounts: newAccountFlags(),
		period:   newDateRangeFlags(unixEpoch, time.Now()),
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
	c.accounts.addFlags(flags)
	c.period.addFlags(flags)

	flags.StringVarP(&c.ledgerDir, "ledger-dir", "", os.Getenv("LEDGER_DIR"), "set ledger directory")

}

func (c *balanceCmd) Execute() error {
	cfg, err := ledger.LoadConfig(filepath.Join(c.ledgerDir, "config.json"))
	if err != nil {
		return err
	}

	var iter ledger.RecordIterator = ledger.NewFilteredIterator(
		ledger.NewLedgerIterator(cfg.Assets, c.ledgerDir),
		c.accounts.accountFilter().Filter,
		c.period.dateRangeFilter().Filter,
	)

	ledger.BalanceReport(iter, cfg.Assets)

	return nil
}
