package main

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	recordTypeInvalid recordType = iota
	recordTypeIncome
	recordTypeExpense
	recordTypeTransfer
)

var (
	ctoi = map[bool]int8{true: 2}
	dtoi = map[bool]int8{true: 1}
)

type recordType int8

type record struct {
	debit      *account
	credit     *account
	recordedAt time.Time
	amount     decimal.Decimal
}

func (r *record) recordType() recordType {
	return recordType(ctoi[r.credit.asset] | dtoi[r.debit.asset])
}

func (r *record) formatAmount() string {
	var sign string

	switch r.recordType() {
	case recordTypeExpense:
		sign = "-"
	case recordTypeIncome:
		sign = "+"
	default:
		sign = "="
	}

	return fmt.Sprintf("%s%s", sign, r.amount.StringFixed(2))
}

func (r *record) String() string {
	return fmt.Sprintf(
		"%q,%q,%s,%s",
		r.debit.name,
		r.credit.name,
		r.amount.String(),
		r.recordedAt.Format(time.RubyDate),
	)
}

func (r *record) toArray() []string {
	return []string{
		r.recordedAt.Format(iso8601Date),
		r.credit.name,
		r.debit.name,
		r.amount.String(),
	}
}
