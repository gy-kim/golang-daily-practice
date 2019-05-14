package stats

import "time"

type Collector struct {
	HostAndPort string
}

func (c *Collector) Count(key string, value int) {
}

func (c *Collector) Duration(key string, start time.Time) {

}
