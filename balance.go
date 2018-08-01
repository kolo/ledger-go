package main

import (
	"fmt"
)

func balanceReport(rd recordReader) {
	balance := report{}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		balance.update(r)
	}

	for _, bi := range balance {
		if bi.account.asset {
			fmt.Printf("%s: %v\n", bi.account.name, bi.total)
		}
	}
}
