package cmd

import (
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/shopspring/decimal"
	"github.com/spf13/pflag"
)

const (
	maxUint64 = ^uint64(0)
	maxInt64  = int64(maxUint64 >> 1)
)

type accountFlags struct {
	credit *string
	debit  *string
}

func newAccountFlags() *accountFlags {
	return &accountFlags{}
}

func (f *accountFlags) addFlags(flags *pflag.FlagSet) {
	f.credit = flags.StringP("credit", "", "", "filter by the credit account")
	f.debit = flags.StringP("debit", "", "", "filter by the debit account")
}

func (f *accountFlags) accountFilter() *ledger.AccountFilter {
	return &ledger.AccountFilter{
		Credit: *f.credit,
		Debit:  *f.debit,
	}
}

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
	flags.Var(f.since, "since", "set a start date for a reporting period")
	flags.Var(f.until, "until", "set an end date for a reporting period")
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
	return time.Time(*f).Format(ledger.ISO8601Date)
}

func (f *dateValue) Set(value string) error {
	t, err := time.Parse(ledger.ISO8601Date, value)
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

type amountFlags struct {
	min *decimalValue
	max *decimalValue
}

func newAmountFlags() *amountFlags {
	return &amountFlags{
		min: newDecimalValue(decimal.Zero),
		max: newDecimalValue(decimal.NewFromInt(maxInt64)),
	}
}

func (f *amountFlags) addFlags(flags *pflag.FlagSet) {
	flags.Var(f.min, "min", "set minimum amount value")
	flags.Var(f.max, "max", "set maximum amount value")
}

func (f *amountFlags) amountFilter() *ledger.AmountFilter {
	return &ledger.AmountFilter{
		Min: f.min.ToDecimal(),
		Max: f.max.ToDecimal(),
	}
}

type decimalValue decimal.Decimal

func newDecimalValue(d decimal.Decimal) *decimalValue {
	p := &d
	return (*decimalValue)(p)
}

func (v *decimalValue) String() string {
	return v.ToDecimal().StringFixed(2)
}

func (v *decimalValue) Set(value string) error {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return err
	}

	*v = decimalValue(d)

	return nil
}

func (v *decimalValue) Type() string {
	return "decimal"
}

func (v *decimalValue) ToDecimal() decimal.Decimal {
	return decimal.Decimal(*v)
}
