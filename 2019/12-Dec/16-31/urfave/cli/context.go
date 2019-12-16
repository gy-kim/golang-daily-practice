package cli

import (
	"context"
	"flag"
)

// Context is a type that is passed through to
// each Handler action in a cli application.
type Context struct {
	context.Context
	App           *App
	Command       *Command
	setFlags      map[string]bool
	flagSet       *flag.FlagSet
	parentContext *Context
}
