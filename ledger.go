package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type ledgerEntry struct {
	creditAccount string
	debitAccount  string
	performedAt   time.Time
	amount        decimal.Decimal
}

func readLedgerEntries(dirs []string) ([]*ledgerEntry, error) {
	entries := []*ledgerEntry{}

	out := []<-chan string{}
	for _, dir := range dirs {
		out = append(out, readLedgerDir(dir))
	}

	var n int64
	for range merge(out...) {
		n++
	}
	fmt.Println(n)

	return entries, nil
}

func readLedgerDir(dir string) <-chan string {
	out := make(chan string)
	go func() {
		err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !fi.IsDir() {
				lines, err := readLedgerFile(p)
				if err != nil {
					return err
				}

				for _, l := range lines {
					out <- l
				}
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
		}

		close(out)
	}()

	return out
}

func readLedgerFile(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := []string{}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for s := range c {
			out <- s
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
