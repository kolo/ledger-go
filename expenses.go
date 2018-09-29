package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// expensesCommand implements "expenses" command.
type expensesCommand struct {
	cmd *cobra.Command
	env *environment

	// filters
	credit *string
	debit  *string
	from   *dateFlag
}

type filterFunc func(*record) *record

func newExpensesCommand(env *environment) *cobra.Command {
	c := &expensesCommand{
		env:  env,
		from: &dateFlag{},
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

	flags.VarP(c.from, "from", "", "set a starting date")

	c.credit = flags.StringP("credit", "", "", "filter by the credit account")
	c.debit = flags.StringP("debit", "", "", "filter by the debit account")
}

func (c *expensesCommand) expenses() {
	expensesReport(c.env.reader(), c.assets(), c.filter())
}

func (c *expensesCommand) assets() []string {
	credit := *c.credit
	if credit == "" {
		return c.env.Accounts
	}

	return []string{credit}
}

func (c *expensesCommand) filter() filterFunc {
	from := c.from.value
	debit := *c.debit

	return func(r *record) *record {
		if r == nil {
			return nil
		}

		if r.recordType() != recordTypeExpense {
			return nil
		}

		if r.recordedAt.Before(from) {
			return nil
		}

		if debit != "" && debit != r.debit.name {
			return nil
		}

		return r
	}
}

func expensesReport(rd recordReader, assets []string, filter filterFunc) {
	expenses := report{}

	for _, asset := range assets {
		expenses[asset] = &reportItem{
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

		if ri, found := expenses[r.credit.name]; found {
			ri.increase(r.amount)
		}
	}

	printExpensesReport(expenses)
}

func printExpensesReport(expenses report) {
	// Sort accounts
	accounts := []string{}
	for k := range expenses {
		accounts = append(accounts, k)
	}
	sort.Strings(accounts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, account := range accounts {
		ri := expenses[account]
		if !ri.total.Equal(decimal.Zero) {
			fmt.Fprintf(w, "%s\t%s\n", ri.account.name, ri.total.StringFixed(2))
		}
	}
	fmt.Fprintln(w, "-----\t")
	fmt.Fprintf(w, "Total\t%s\n", expenses.total().StringFixed(2))

	w.Flush()
}
