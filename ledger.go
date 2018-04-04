package main

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type recordType int8

const (
	recordTypeInvalid recordType = iota
	recordTypeIncome
	recordTypeExpense
	recordTypeTransfer
)

var (
	ctoi = map[bool]int8{true: 2}
	dtoi = map[bool]int8{true: 1}

	minusOne = decimal.NewFromFloat(-1.0)
)

// balance prints balance of each account in the ledger.
func balance(assets []string, records []string) {
	assetsBalances := map[string]decimal.Decimal{}
	expensesBalances := map[string]decimal.Decimal{}
	incomesBalances := map[string]decimal.Decimal{}

	assetsMap := map[string]bool{}
	for _, asset := range assets {
		assetsMap[asset] = true
	}

	for _, record := range records {
		values := strings.Split(record, ",")
		creditAccount := values[1]
		debitAccount := values[2]
		amount, _ := decimal.NewFromString(values[3])
		// calculate record type
		_, isCreditAnAsset := assetsMap[creditAccount]
		_, isDebitAnAsset := assetsMap[debitAccount]
		rt := recordType(ctoi[isCreditAnAsset] | dtoi[isDebitAnAsset])

		switch rt {
		case recordTypeExpense:
			updateBalance(creditAccount, assetsBalances, amount.Mul(minusOne))
			updateBalance(debitAccount, expensesBalances, amount)
		case recordTypeIncome:
			updateBalance(creditAccount, incomesBalances, amount.Mul(minusOne))
			updateBalance(debitAccount, assetsBalances, amount)
		case recordTypeTransfer:
			updateBalance(creditAccount, assetsBalances, amount.Mul(minusOne))
			updateBalance(debitAccount, assetsBalances, amount)
		}
	}

	fmt.Println("Assets:")
	for _, asset := range assets {
		fmt.Printf("  %s: %s\n", asset, assetsBalances[asset].String())
	}

	fmt.Println("Income:")
	for account, balance := range incomesBalances {
		fmt.Printf("  %s: %s\n", account, balance.String())
	}

	fmt.Println("Expenses:")
	for account, balance := range expensesBalances {
		fmt.Printf("  %s: %s\n", account, balance.String())
	}
}

func updateBalance(account string, balances map[string]decimal.Decimal, amount decimal.Decimal) {
	if _, registered := balances[account]; !registered {
		balances[account] = decimal.Zero
	}

	balances[account] = balances[account].Add(amount)
}
