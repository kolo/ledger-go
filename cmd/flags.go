package cmd

import (
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/pflag"
)

const iso8601Date = "2006-01-02"

type dateRangeFlags struct {
	since *dateValue
	until *dateValue
}

func newDateRangeFlags(since, until time.Time) *dateRangeFlags {
	return &dateRangeFlags{
		since: newDateValue(since),
		until: newDateValue(until),
	}
}

func (f *dateRangeFlags) addFlags(flags *pflag.FlagSet) {
	flags.VarP(f.since, "since", "", "set a start date for a reporting period")
	flags.VarP(f.until, "until", "", "set an end date for a reporting period")
}

func (f *dateRangeFlags) dateRangeFilter() *ledger.DateRangeFilter {
	return &ledger.DateRangeFilter{
		Since: f.since.ToTime(),
		Until: f.until.ToTime(),
	}
}

type dateValue time.Time

func newDateValue(t time.Time) *dateValue {
	p := &t
	return (*dateValue)(p)
}

func (f *dateValue) String() string {
	return time.Time(*f).Format(iso8601Date)
}

func (f *dateValue) Set(value string) error {
	t, err := time.Parse(iso8601Date, value)
	if err != nil {
		return err
	}

	*f = dateValue(t)

	return nil
}

func (f *dateValue) Type() string {
	return "date"
}

func (f *dateValue) ToTime() time.Time {
	return time.Time(*f)
}
