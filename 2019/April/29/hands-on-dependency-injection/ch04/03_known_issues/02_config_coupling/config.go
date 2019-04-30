package config_coupling

import (
	"github.com/gy-kim/golang-daily-practice/2019/April/29/hands-on-dependency-injection/ch04/03_known_issues/02_config_coupling/currency"
)

type Config struct {
	DefaultCurrency currency.Currency `json:"default_currency"`
}
