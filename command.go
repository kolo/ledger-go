package main

import (
	"time"

	"github.com/spf13/cobra"
)

const iso8601Date = "2006-01-02"

type command interface {
	Cmd() *cobra.Command
}

// dateFlag represents a date flag value encoded in iso8601 date format.
type dateFlag struct {
	value time.Time
}

func (f *dateFlag) String() string {
	return f.value.Format(iso8601Date)
}

// Set ...
func (f *dateFlag) Set(value string) error {
	var err error

	f.value, err = time.Parse(iso8601Date, value)
	if err != nil {
		return err
	}

	return nil
}

// Type ...
func (f *dateFlag) Type() string {
	return "date"
}
