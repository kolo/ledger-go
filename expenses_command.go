package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// expensesCommand implements "expenses" command.
type expensesCommand struct {
	cmd *cobra.Command
	env *environment

	// filters
	accounts *accountFilter
	period   *dateRangeFilter

	// grouping
	weekly *bool
}

type filterFunc func(*record) *record

func newExpensesCommand(env *environment) *cobra.Command {
	c := &expensesCommand{
		env:      env,
		accounts: newAccountFilter(),
		period:   newDateRangeFilter(),
	}

	c.cmd = &cobra.Command{
		Use: "expenses [OPTIONS]",
		Run: func(*cobra.Command, []string) {
			c.expenses()
		},
	}
	c.addFlags()

	return c.Cmd()
}

// Cmd initializes an instance of cobra.Command.
func (c *expensesCommand) Cmd() *cobra.Command {
	return c.cmd
}

func (c *expensesCommand) addFlags() {
	flags := c.cmd.Flags()

	c.accounts.addFlags(flags)
	c.period.addFlags(flags)

	c.weekly = flags.BoolP("weekly", "", false, "group expenses by week")
}

func (c *expensesCommand) expenses() {
	reader := newFilteredReader(c.env.reader(), c.filter())

	if *c.weekly {
		weeklyExpensesReport(reader, c.assets())
		return
	}

	expensesReport(reader, c.assets())
}

func (c *expensesCommand) assets() []string {
	return c.env.Accounts
}

func (c *expensesCommand) filter() filterFunc {
	return func(r *record) *record {
		if r == nil {
			return nil
		}

		if r.recordType() != recordTypeExpense {
			return nil
		}

		r = c.accounts.filter(r)
		if r == nil {
			return nil
		}

		r = c.period.filter(r)
		if r == nil {
			return nil
		}

		return r
	}
}

func expensesReport(rd recordReader, assets []string) {
	credits, debits := map[string]struct{}{}, map[string]struct{}{}
	expenses := map[[2]string]decimal.Decimal{}

	update := func(r *record) {
		key := [2]string{r.credit.name, r.debit.name}
		if _, found := expenses[key]; !found {
			expenses[key] = decimal.Zero
		}

		if _, found := credits[r.credit.name]; !found {
			credits[r.credit.name] = struct{}{}
		}

		if _, found := debits[r.debit.name]; !found {
			debits[r.debit.name] = struct{}{}
		}

		expenses[key] = expenses[key].Add(r.amount)
	}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		update(r)
	}

	// printing
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	columns := []string{}
	totals := map[string]decimal.Decimal{}
	for credit := range credits {
		columns = append(columns, credit)
		totals[credit] = decimal.Zero
	}
	sort.Strings(columns)

	rows := []string{}
	for debit := range debits {
		rows = append(rows, debit)
	}
	sort.Strings(rows)

	// header
	fmt.Fprintf(w, "\t%v\t=\n", strings.Join(columns, "\t"))

	// body
	for _, debit := range rows {
		v := []string{debit}
		t := decimal.Zero

		for _, credit := range columns {
			key := [2]string{credit, debit}
			spend := expenses[key]
			if spend.Equals(decimal.Zero) {
				v = append(v, "-")
			} else {
				v = append(v, spend.StringFixed(2))
				t = t.Add(spend)
				totals[credit] = totals[credit].Add(spend)
			}
		}

		fmt.Fprintf(w, "%s\t%s\n", strings.Join(v, "\t"), t.StringFixed(2))
	}

	// footer
	fmt.Fprintf(w, "=")
	for _, column := range columns {
		fmt.Fprintf(w, "\t%s", totals[column].StringFixed(2))
	}
	fmt.Fprintf(w, "\n")

	w.Flush()
}
