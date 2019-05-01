package main

import (
	"time"

	"github.com/spf13/pflag"
)

type dateRangeFilter struct {
	since *dateFlag
	until *dateFlag
}

func (f *dateRangeFilter) addFlags(flags *pflag.FlagSet) {
	flags.VarP(f.since, "since", "", "set a start date for a reporting period")
	flags.VarP(f.until, "until", "", "set an end date for a reporting period")
}

func (f *dateRangeFilter) filter(r *record) *record {
	if r.recordedAt.Before(f.start()) || r.recordedAt.After(f.end()) {
		return nil
	}

	return r
}

func (f *dateRangeFilter) start() time.Time {
	return f.since.value
}

func (f *dateRangeFilter) end() time.Time {
	return f.until.value
}

func newDateRangeFilter() *dateRangeFilter {
	return &dateRangeFilter{
		since: &dateFlag{},
		until: &dateFlag{value: time.Now()},
	}
}
