package ledger

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	yearRx  = regexp.MustCompile("^\\d{4}$")
	monthRx = regexp.MustCompile("^[0-1][0-9]$")
	dayRx   = regexp.MustCompile("^[0-3][0-9]$")
)

type RecordIterator interface {
	Next() *Record
}

type LedgerIterator struct {
	assets  map[string]*Account
	records []string

	cur int
}

func NewLedgerIterator(assets []*Account, ledgerDir string) *LedgerIterator {
	assetsMap := map[string]*Account{}
	for _, asset := range assets {
		assetsMap[asset.Name] = asset
	}

	return &LedgerIterator{
		assets:  assetsMap,
		records: readLedgerDir(ledgerDir),
	}
}

func (i *LedgerIterator) Next() *Record {
	if i.cur == len(i.records) {
		// end of list
		return nil
	}

	values := strings.Split(i.records[i.cur], ",")
	i.cur = i.cur + 1

	amount, _ := decimal.NewFromString(values[3])
	date, _ := time.Parse("2006-01-02", values[0])
	return &Record{
		Credit: i.newAccount(values[1]),
		Debit:  i.newAccount(values[2]),
		Date:   date,
		Amount: amount,
	}
}

func (i *LedgerIterator) newAccount(name string) *Account {
	account, found := i.assets[name]
	if !found {
		account = NewAccount(name, false)
	}

	return account
}

func readLedgerDir(ledgerDir string) []string {
	paths := lookupLedgerFiles(ledgerDir, yearRx, monthRx, dayRx)

	records := []string{}
	for _, path := range paths {
		date := pathToDate(ledgerDir, path)
		records = append(records, readLedgerFile(path, date)...)
	}
	return records
}

func lookupLedgerFiles(root string, patterns ...*regexp.Regexp) []string {
	matches := []string{root}
	for _, pattern := range patterns {
		n := len(matches)
		for _, path := range matches {
			entries, err := ioutil.ReadDir(path)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if pattern.MatchString(entry.Name()) {
					matches = append(matches, filepath.Join(path, entry.Name()))
				}
			}
		}
		matches = matches[n:]
	}

	return matches
}

func pathToDate(basepath string, path string) string {
	rel, _ := filepath.Rel(basepath, path)
	return strings.Replace(rel, string(filepath.Separator), "-", -1)
}

func readLedgerFile(path string, date string) []string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}
	}

	sc := bufio.NewScanner(bytes.NewReader(b))

	records := []string{}
	for sc.Scan() {
		r := strings.Join([]string{date, sc.Text()}, ",")
		records = append(records, r)
	}

	return records
}
