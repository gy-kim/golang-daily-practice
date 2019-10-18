package logrus

import (
	"runtime"
	"sync"
	"time"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	ForceColors               bool
	DisableColors             bool
	EnvironmentOverrideColors bool
	DisableTimestamp          bool
	FullTimestamp             bool
	TimestampFormat           string
	DisableSorting            bool
	SortingFunc               func([]string)
	DisableLevelTruncation    bool
	QuoteEmptyFields          bool
	isTerminal                bool
	FieldMap                  FieldMap
	CallerPrettyfier          func(*runtime.Frame) (function string, file string)
	ternimalInitOnce          sync.Once
}
