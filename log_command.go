package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type logCommand struct {
	cmd *cobra.Command
	env *environment

	period *dateRangeFilter
}

func newLogCommand(env *environment) *cobra.Command {
	c := &logCommand{
		env:    env,
		period: newDateRangeFilter(),
	}

	c.cmd = &cobra.Command{
		Use: "log [OPTIONS]",
		Run: func(*cobra.Command, []string) {
			c.log()
		},
	}
	c.addFlags()

	return c.cmd
}

func (c *logCommand) addFlags() {
	flags := c.cmd.Flags()
	c.period.addFlags(flags)
}

func (c *logCommand) log() {
	rd := newFilteredReader(c.env.reader(), c.period.filter)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		date := r.recordedAt.Format(iso8601Date)
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\n",
			date,
			r.credit.name,
			r.debit.name,
			c.formatAmount(r),
		)

	}

	w.Flush()
}

func (c *logCommand) formatAmount(r *record) string {
	var sign string
	switch r.recordType() {
	case recordTypeExpense:
		sign = "-"
	case recordTypeIncome:
		sign = "+"
	default:
		sign = "="
	}

	return fmt.Sprintf("%s%s", sign, r.amount.StringFixed(2))
}
