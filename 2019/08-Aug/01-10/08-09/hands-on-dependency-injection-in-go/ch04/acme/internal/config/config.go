package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/01-10/08-09/hands-on-dependency-injection-in-go/ch04/acme/internal/logging"
)

const DefaultEnvVar = "ACME_CONFIG"

var App *Config

type Config struct {
	DSN string

	Address string

	BasePrice float64

	ExchangeRateBaseURL string

	ExchangeRateAPIKey string
}

func init() {
	filename, found := os.LookupEnv(DefaultEnvVar)
	if !found {
		logging.L.Error("failed to locate file specified by %s", DefaultEnvVar)
		return
	}
	_ = load(filename)
}

func load(filename string) error {
	App = &Config{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logging.L.Error("failed to read config file. err: %s", err)
		return err
	}

	err = json.Unmarshal(bytes, App)
	if err != nil {
		logging.L.Error("failed to parse config file. err: %s", err)
		return err
	}

	return nil
}
