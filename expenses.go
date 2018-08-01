package main

import (
	"fmt"
	"time"

	"github.com/kolo/ledger-go/datetime"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// expensesCommand implements "expenses" command.
type expensesCommand struct {
	cmd *cobra.Command

	from *dateFlag
}

func newExpensesCommand() *cobra.Command {
	c := &expensesCommand{
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
}

func (c *expensesCommand) expenses() {
	fmt.Println(c.from.value.Format(iso8601Date))
}

type expenses map[string]*reportItem

func (e expenses) update(r *record) {
	if r.recordType() != recordTypeExpense {
		return
	}

	if bi, found := e[r.credit.name]; found {
		bi.increase(r.amount)
	}
}

func (e expenses) total() decimal.Decimal {
	total := decimal.Zero
	for _, bi := range e {
		total = total.Add(bi.total)
	}

	return total
}

func expensesReport(rd recordReader, assets []string) {
	expenses := expenses{}
	for _, asset := range assets {
		expenses[asset] = &reportItem{
			account: &account{
				name:  asset,
				asset: true,
			},
			total: decimal.Zero,
		}
	}

	t := datetime.BeginningOfWeek(time.Now())
	for {
		r := rd.Next()
		if r == nil {
			break
		}

		if r.recordedAt.After(t) {
			expenses.update(r)
		}
	}

	for _, bi := range expenses {
		fmt.Printf("%5s: %6s\n", bi.account.name, bi.total.StringFixed(2))
	}
	fmt.Printf("Total: %v\n", expenses.total())
}
