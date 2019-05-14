package config_injection

import (
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/logging"
	"github.com/gy-kim/golang-daily-practice/2019/05-May/11-20/14/hands-on-dependency-injection-in-go/ch08/02_advantages/stats"
)

func NewMyObject(cfg Config) *MyObject {
	return &MyObject{
		cfg: cfg,
	}
}

type Config interface {
	Logger() *logging.Logger
	Stats() *stats.Collector
}

type MyObject struct {
	cfg Config
}

func (m *MyObject) Do() (interface{}, error) {
	m.cfg.Logger().Error("not implemented")
	return struct{}{}, nil
}
