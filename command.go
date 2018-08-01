package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const iso8601Date = "2006-01-02"

// Command ...
type Command interface {
	Cmd() *cobra.Command
}

// ExpensesCommand implements "expenses" command.
type ExpensesCommand struct {
	cmd *cobra.Command

	from *DateFlag
}

// Cmd initializes an instance of cobra.Command.
func (c *ExpensesCommand) Cmd() *cobra.Command {
	c.cmd = &cobra.Command{
		Use: "expenses [OPTIONS]",
		Run: func(*cobra.Command, []string) {
			c.balance()
		},
	}
	c.addFlags()

	return c.cmd
}

func (c *ExpensesCommand) addFlags() {
	flags := c.cmd.Flags()

	c.from = &DateFlag{}
	flags.VarP(c.from, "from", "", "set a starting date")
}

func (c *ExpensesCommand) balance() {
	fmt.Println(c.from.value.Format(iso8601Date))
}

// DateFlag represents a date flag value encoded in iso8601 date format.
type DateFlag struct {
	value time.Time
}

func (f *DateFlag) String() string {
	return f.value.Format(iso8601Date)
}

// Set ...
func (f *DateFlag) Set(value string) error {
	var err error

	f.value, err = time.Parse(iso8601Date, value)
	if err != nil {
		return err
	}

	return nil
}

// Type ...
func (f *DateFlag) Type() string {
	return "date"
}
