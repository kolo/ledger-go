package ledger

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type configValues struct {
	Accounts []string `json:"accounts"`
}

type Config struct {
	Assets []*Account
}

func LoadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("Can't read %s", path)
	}

	v := configValues{Accounts: []string{}}
	if err = json.Unmarshal(b, &v); err != nil {
		return nil, errors.Errorf("%s is not valid json file", path)
	}

	return newConfig(v), nil
}

func newConfig(v configValues) *Config {
	cfg := &Config{}
	for _, accountName := range v.Accounts {
		cfg.Assets = append(cfg.Assets, NewAccount(accountName, true))
	}

	return cfg
}
