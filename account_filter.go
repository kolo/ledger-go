package main

import "github.com/spf13/pflag"

type accountFilter struct {
	credit *string
	debit  *string
}

func (f *accountFilter) addFlags(flags *pflag.FlagSet) {
	f.credit = flags.StringP("credit", "", "", "filter by the credit account")
	f.debit = flags.StringP("debit", "", "", "filter by the debit account")
}

func (f *accountFilter) filter(r *record) *record {
	if *f.debit != "" {
		if r.debit.name != *f.debit {
			return nil
		}
	}

	if *f.credit != "" {
		if r.credit.name != *f.credit {
			return nil
		}
	}

	return r
}

func newAccountFilter() *accountFilter {
	return &accountFilter{}
}
