package cmd

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/pflag"
)

var unixEpoch = time.Unix(0, 0)

type baseCmd struct {
	accounts *accountFlags
	amount   *amountFlags
	period   *dateRangeFlags

	ledgerDir string

	run func(*ledger.Config, ledger.RecordIterator, io.Writer) error
}

func newBaseCmd() *baseCmd {
	return &baseCmd{
		accounts: newAccountFlags(),
		amount:   newAmountFlags(),
		period:   newDateRangeFlags(unixEpoch, time.Now()),
	}
}

func (c *baseCmd) addFlags(flags *pflag.FlagSet) {
	c.accounts.addFlags(flags)
	c.amount.addFlags(flags)
	c.period.addFlags(flags)

	flags.StringVar(&c.ledgerDir, "ledger-dir", os.Getenv("LEDGER_DIR"), "set ledger directory")

}

func (c *baseCmd) Run(stdout io.Writer) error {
	cfg, err := ledger.LoadConfig(filepath.Join(c.ledgerDir, "config.json"))
	if err != nil {
		return err
	}

	var iter ledger.RecordIterator = ledger.NewFilteredIterator(
		ledger.NewLedgerIterator(cfg.Assets, c.ledgerDir),
		c.accounts.accountFilter().Filter,
		c.amount.amountFilter().Filter,
		c.period.dateRangeFilter().Filter,
	)

	return c.run(cfg, iter, stdout)
}
