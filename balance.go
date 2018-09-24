package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type balanceCommand struct {
	cmd *cobra.Command
	env *environment
}

func newBalanceCommand(env *environment) *cobra.Command {
	c := &balanceCommand{env: env}
	c.cmd = &cobra.Command{
		Use: "balance",
		Run: func(*cobra.Command, []string) {
			c.balance()
		},
	}

	return c.Cmd()
}

func (c *balanceCommand) Cmd() *cobra.Command {
	return c.cmd
}

func (c *balanceCommand) balance() {
	balanceReport(c.env.reader(), c.env.Accounts)
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
