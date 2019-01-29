package main

import (
	"encoding/csv"
	"os"

	"github.com/spf13/cobra"
)

type exportCommand struct {
	cmd *cobra.Command
	env *environment

	period *dateRangeFilter
}

func newExportCommand(env *environment) *cobra.Command {
	c := &exportCommand{
		env:    env,
		period: newDateRangeFilter(),
	}

	c.cmd = &cobra.Command{
		Use: "export [OPTIONS]",
		RunE: func(*cobra.Command, []string) error {
			return c.export()
		},
	}
	c.addFlags()

	return c.cmd
}

func (c *exportCommand) addFlags() {
	flags := c.cmd.Flags()
	c.period.addFlags(flags)
}

func (c *exportCommand) export() error {
	rd := newFilteredReader(c.env.reader(), c.period.filter)

	f, err := os.OpenFile("ledger.csv", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for {
		r := rd.Next()
		if r == nil {
			break
		}

		if err = w.Write(r.toArray()); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}
