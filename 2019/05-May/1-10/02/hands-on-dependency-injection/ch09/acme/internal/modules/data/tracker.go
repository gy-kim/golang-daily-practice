package data

import (
	"time"

	"github.com/gy-kim/golang-daily-practice/2019/05-May/1-10/02/hands-on-dependency-injection/ch09/acme/internal/logging"
)

// QueryTracker is an interface to track query timing
type QueryTracker interface {
	// Track will record/out the time a query took by calculating time.Now().Sub(start)
	Track(key string, start time.Time)
}

// NO-OP implementation of QueryTracker
type noopTracker struct{}

// Track implements QueryTracker
func (_ *noopTracker) Track(_ string, _ time.Time) {
	// intentionally does nothing
}

// NewLogTracker returns a Tracker that outputs tracking data to og
func NewLogTracker(logger logging.Logger) *LogTracker {
	return &LogTracker{
		logger: logger,
	}
}

// LogTracker implements QueryTracker and outputs to the supplied logger
type LogTracker struct {
	logger logging.Logger
}

// Track implements QueryTracker
func (l *LogTracker) Track(key string, start time.Time) {
	l.logger.Info("[%s] Timing: %s\n", key, time.Now().Sub(start))
}
