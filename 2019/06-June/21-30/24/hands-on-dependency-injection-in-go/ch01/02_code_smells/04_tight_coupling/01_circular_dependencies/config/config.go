package config

import (
	"errors"

	"github.com/gy-kim/golang-daily-practice/2019/06-June/21-30/24/hands-on-dependency-injection-in-go/ch01/02_code_smells/04_tight_coupling/01_circular_dependencies/payment"
)

// Config defines the JSON format of the config file
type Config struct {
	// Address is the host and port to bind to.
	// Default 0.0.0.0:8080
	Address string

	// DefaultCurrency is the default currency of the system
	DefaultCurrency payment.Currency
}

// Load will load the JSON config from the file supplied
func Load(fulename string) (*Config, error) {
	// TODO: load currency from file
	return nil, errors.New("not implemented yet")
}
