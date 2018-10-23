package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

type registerCommand struct {
	cmd *cobra.Command
	env *environment

	period *dateRangeFilter
}

func newRegisterCommand(env *environment) *cobra.Command {
	c := &registerCommand{
		env:    env,
		period: newDateRangeFilter(),
	}

	c.cmd = &cobra.Command{
		Use: "register [OPTIONS]",
		Run: func(*cobra.Command, []string) {
			c.register()
		},
	}
	c.addFlags()

	return c.cmd
}

func (c *registerCommand) addFlags() {
	flags := c.cmd.Flags()
	c.period.addFlags(flags)
}

func (c *registerCommand) register() {
	reader := newFilteredReader(c.env.reader(), c.period.filter)
	total := decimal.Zero

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for {
		r := reader.Next()
		if r == nil {
			break
		}

		switch r.recordType() {
		case recordTypeExpense:
			total = total.Sub(r.amount)
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				r.recordedAt.Format("02-Jan-06"),
				r.debit.name,
				fmt.Sprintf("-%s", r.amount.StringFixed(2)),
				total.StringFixed(2),
			)
		case recordTypeIncome:
			total = total.Add(r.amount)
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\n",
				r.recordedAt.Format("02-Jan-06"),
				r.debit.name,
				fmt.Sprintf("+%s", r.amount.StringFixed(2)),
				total.StringFixed(2),
			)
		}
	}

	w.Flush()
}
