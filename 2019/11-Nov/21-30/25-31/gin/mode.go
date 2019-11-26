package gin

import (
	"io"
	"os"
)

// EnvGinMode indicates environment name for gin mode.
const EnvGinMode = "GIN_MODE"

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

// DefaultWriter is the default io.Writer used by Gin for debug output and
// middleware output like Logger() or Recovery().
var DefaultWriter io.Writer = os.Stderr

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

var ginMode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv(EnvGinMode)
	SetMode(mode)
}

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		ginMode = debugCode
	}
}

// DisableBindValidation closes the default validator.
func DisableBindValidation() {
	binding.Validator = nil
}
