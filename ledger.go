package main

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type account struct {
	name  string
	asset bool
}

type recordReader interface {
	Next() *record
}

type simpleReader struct {
	accounts map[string]*account
	records  []string

	cur int
}

func (r *simpleReader) Next() *record {
	if r.cur == len(r.records) {
		// end of list
		return nil
	}

	values := strings.Split(r.records[r.cur], ",")
	r.cur = r.cur + 1

	amount, _ := decimal.NewFromString(values[3])
	recordedAt, _ := time.Parse("2006-01-02", values[0])
	return &record{
		credit:     r.account(values[1]),
		debit:      r.account(values[2]),
		amount:     amount,
		recordedAt: recordedAt,
	}
}

func (r *simpleReader) account(name string) *account {
	ac, found := r.accounts[name]
	if !found {
		ac = &account{
			name:  name,
			asset: false,
		}
	}

	return ac
}

func newSimpleReader(assets []string, records []string) *simpleReader {
	accounts := map[string]*account{}
	for _, asset := range assets {
		accounts[asset] = &account{
			name:  asset,
			asset: true,
		}
	}

	return &simpleReader{
		accounts: accounts,
		records:  records,
	}
}
