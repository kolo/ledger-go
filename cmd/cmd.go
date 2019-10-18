package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/pflag"
)

type baseCmd struct {
	accounts *accountFlags
	period   *dateRangeFlags

	ledgerDir string

	run func(*ledger.Config, ledger.RecordIterator) error
}

func newBaseCmd() *baseCmd {
	return &baseCmd{
		accounts: newAccountFlags(),
		period:   newDateRangeFlags(unixEpoch, time.Now()),
	}
}

func (c *baseCmd) addFlags(flags *pflag.FlagSet) {
	c.accounts.addFlags(flags)
	c.period.addFlags(flags)

	flags.StringVarP(&c.ledgerDir, "ledger-dir", "", os.Getenv("LEDGER_DIR"), "set ledger directory")

}

func (c *baseCmd) Run() error {
	cfg, err := ledger.LoadConfig(filepath.Join(c.ledgerDir, "config.json"))
	if err != nil {
		return err
	}

	var iter ledger.RecordIterator = ledger.NewFilteredIterator(
		ledger.NewLedgerIterator(cfg.Assets, c.ledgerDir),
		c.accounts.accountFilter().Filter,
		c.period.dateRangeFilter().Filter,
	)

	return c.run(cfg, iter)
}
