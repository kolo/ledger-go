package ledger

import (
	"time"

	"github.com/shopspring/decimal"
)

type RecordFilter func(*Record) *Record

type AccountFilter struct {
	Credit string
	Debit  string
}

func (f *AccountFilter) Filter(r *Record) *Record {
	if f.Credit != "" && r.Credit.Name != f.Credit {
		return nil
	}

	if f.Debit != "" && r.Debit.Name != f.Debit {
		return nil
	}

	return r
}

type DateRangeFilter struct {
	Since time.Time
	Until time.Time
}

func (f *DateRangeFilter) Filter(r *Record) *Record {
	if r.Date.Before(f.Since) || r.Date.After(f.Until) {
		return nil
	}

	return r
}

type AmountFilter struct {
	Min decimal.Decimal
	Max decimal.Decimal
}

func (f *AmountFilter) Filter(r *Record) *Record {
	if r.Amount.LessThanOrEqual(f.Max) && r.Amount.GreaterThanOrEqual(f.Min) {
		return r
	}

	return nil
}
