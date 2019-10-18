package ledger

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	recordTypeInvalid RecordType = iota
	recordTypeIncome
	recordTypeExpense
	recordTypeTransfer
)

var (
	ctoi = map[bool]int8{true: 2}
	dtoi = map[bool]int8{true: 1}
)

type RecordType int8

type Record struct {
	Debit  *Account
	Credit *Account
	Date   time.Time
	Amount decimal.Decimal
}

func (r *Record) RecordType() RecordType {
	return RecordType(ctoi[r.Credit.Asset] | dtoi[r.Debit.Asset])
}

func (r *Record) formatAmount() string {
	var sign string

	switch r.RecordType() {
	case recordTypeExpense:
		sign = "-"
	case recordTypeIncome:
		sign = "+"
	default:
		sign = "="
	}

	return fmt.Sprintf("%s%s", sign, r.Amount.StringFixed(2))
}
