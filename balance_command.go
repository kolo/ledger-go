package main

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type balanceCommand struct {
	cmd *cobra.Command
	env *environment
}

func newBalanceCommand(env *environment) *cobra.Command {
	c := &balanceCommand{env: env}

	period := newDateRangeFilter()

	c.cmd = &cobra.Command{
		Use: "balance",
		Run: func(*cobra.Command, []string) {
			c.balance(period.filter)
		},
	}

	period.addFlags(c.cmd.Flags())

	return c.cmd
}

func (c *balanceCommand) balance(filter filterFunc) {
	balanceReport(c.env.reader(), c.env.Accounts, filter)
}

func balanceReport(rd recordReader, assets []string, filter filterFunc) {
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

		r = filter(r)
		if r == nil {
			continue
		}

		if ri, ok := balance[r.credit.name]; ok {
			ri.decrease(r.amount)
		}
		if ri, ok := balance[r.debit.name]; ok {
			ri.increase(r.amount)
		}
	}

	printReport(balance)
}
