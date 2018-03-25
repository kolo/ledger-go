package main

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	minusOne = decimal.NewFromFloat(-1.0)
)

func total(accounts []string, records []string) {
	totals := map[string]decimal.Decimal{}
	for _, account := range accounts {
		totals[account] = decimal.Zero
	}

	for _, r := range records {
		values := strings.Split(r, ",")

		creditAccount := values[1]
		debitAccount := values[2]
		amount, err := decimal.NewFromString(values[3])
		if err != nil {
			fmt.Printf("Can't decode: %s\n", r)
		}

		if _, ok := totals[creditAccount]; ok {
			totals[creditAccount] = totals[creditAccount].Sub(amount)
		}

		if _, ok := totals[debitAccount]; ok {
			totals[debitAccount] = totals[debitAccount].Add(amount)
		}
	}

	for account, total := range totals {
		fmt.Printf("%s: %s\n", account, total.String())
	}
}
