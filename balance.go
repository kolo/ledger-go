package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type balanceCommand struct {
	cmd *cobra.Command
}

func newBalanceCommand() *cobra.Command {
	c := &balanceCommand{}
	c.cmd = &cobra.Command{
		Use: "balance",
		RunE: func(*cobra.Command, []string) error {
			return c.balance()
		},
	}

	return c.Cmd()
}

func (c *balanceCommand) Cmd() *cobra.Command {
	return c.cmd
}

func (c *balanceCommand) balance() error {
	ledgerDir := os.Getenv(ledgerDirEnvKey)
	if ledgerDir == "" {
		return errors.Errorf("missing environment variable - %s", ledgerDirEnvKey)
	}

	config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
	if err != nil {
		return errors.Wrap(err, "can't read configuration")
	}

	reader := newSimpleReader(config.Accounts, read(ledgerDir))
	balanceReport(reader, config.Accounts)

	return nil
}

func balanceReport(rd recordReader, assets []string) {
	balance := report{}

	for _, asset := range assets {
		balance[asset] = &reportItem{
			account: &account{
				name:  asset,
				asset: true,
			},
			total: decimal.Zero,
		}
	}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		if ri, ok := balance[r.credit.name]; ok {
			ri.decrease(r.amount)
		}
		if ri, ok := balance[r.debit.name]; ok {
			ri.increase(r.amount)
		}
	}

	for _, ri := range balance {
		fmt.Printf("%s: %v\n", ri.account.name, ri.total)
	}
}
