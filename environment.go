package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type environment struct {
	ledgerDir string
	*userConfig
}

func loadEnvironment() (*environment, error) {
	ledgerDir := os.Getenv(ledgerDirEnvKey)
	if ledgerDir == "" {
		return nil, errors.Errorf("missing environment variable - %s", ledgerDirEnvKey)
	}

	config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
	if err != nil {
		return nil, errors.Wrap(err, "can't read configuration")
	}

	return &environment{
		ledgerDir:  ledgerDir,
		userConfig: config,
	}, nil
}

func (e *environment) reader() recordReader {
	return newSimpleReader(e.Accounts, read(e.ledgerDir))
}
