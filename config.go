package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

type userConfig struct {
	Accounts []string `json:"accounts"`
}

func loadUserConfig(filepath string) (*userConfig, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't read %s", filepath))
	}

	c := &userConfig{Accounts: []string{}}
	if err = json.Unmarshal(b, c); err != nil {
		return nil, errors.New(fmt.Sprintf("%s is not valid json file", filepath))
	}

	return c, nil
}
