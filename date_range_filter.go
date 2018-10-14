package main

import (
	"time"

	"github.com/spf13/pflag"
)

type dateRangeFilter struct {
	from *dateFlag
	to   *dateFlag
}

func (f *dateRangeFilter) addFlags(flags *pflag.FlagSet) {
	flags.VarP(f.from, "from", "", "set a start date for a reporting period")
	flags.VarP(f.to, "to", "", "set an end date for a reporting period")
}

func (f *dateRangeFilter) filter(r *record) *record {
	if r.recordedAt.Before(f.start()) || r.recordedAt.After(f.end()) {
		return nil
	}

	return r
}

func (f *dateRangeFilter) start() time.Time {
	return f.from.value
}

func (f *dateRangeFilter) end() time.Time {
	return f.to.value
}
