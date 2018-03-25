package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	yearRx  = regexp.MustCompile("^\\d{4}$")
	monthRx = regexp.MustCompile("^[0-1][0-9]$")
	dayRx   = regexp.MustCompile("^[0-3][0-9]$")
)

func read(ledgerDir string) []string {
	paths := lookup(ledgerDir, yearRx, monthRx, dayRx)

	records := []string{}
	for _, path := range paths {
		date := pathToDate(ledgerDir, path)
		records = append(records, readLedgerFile(path, date)...)
	}
	return records
}

func lookup(root string, patterns ...*regexp.Regexp) []string {
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
