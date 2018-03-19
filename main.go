package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

const (
	ledgerDirEnvKey = "LEDGER_DIR"
)

var (
	dataDirRx = regexp.MustCompile("^\\d{4}$")
)

func main() {
	ledgerDir := os.Getenv(ledgerDirEnvKey)
	if ledgerDir == "" {
		exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
	}

	dataDirectories, err := getDataDirectories(ledgerDir)
	if err != nil {
		exitWithErr(err)
	}

	for _, dir := range dataDirectories {
		fmt.Println(dir)
	}

	ledgerEntries, err := readLedgerEntries(dataDirectories)
	if err != nil {
		exitWithErr(err)
	}

	fmt.Println(len(ledgerEntries))
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

func getDataDirectories(baseDir string) ([]string, error) {
	dirs := []string{}

	entries, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't read content of %q", baseDir))
	}

	for _, entry := range entries {
		if entry.IsDir() && dataDirRx.MatchString(entry.Name()) {
			dirs = append(dirs, filepath.Join(baseDir, entry.Name()))
		}
	}

	return dirs, nil
}
