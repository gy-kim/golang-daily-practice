package cli

import (
	"io"
	"time"
)

// App is the main structure of a cli application. It is recommended that
// an app be created with the cli.NewApp() function
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string

	// Full name of command for help, defaults to Name
	HelpName string

	// Description of the program
	Usage string

	// Text to override the USAGE section of help.
	UsageText string

	// Description of the program argument format.
	ArgsUsage string

	// Version of the program
	Version string

	// Description of the program
	Description string

	// List of commands to execute
	Commands []*Command

	// List of flags to parse
	Flags []Flag

	EnableBashCompletion bool

	HideHelp bool

	HideVersion bool

	categories CommandCategories

	BashComplete BashCompleteFunc

	Before BeforeFunc

	After AfterFunc

	CommandNotFound CommandNotFoundFunc

	OnUsageError OnUsageErrorFunc

	Compiled time.Time

	Authors []*Author

	Copyright string

	Writer io.Writer

	ErrWriter io.Writer

	ExitErrHandler ExitErrHandlerFunc

	Metadata map[string]interface{}

	ExtraInfo func() map[string]string

	CustomAppHelpTemplate string

	UseShortOptionHandling bool

	didSetup bool
}

type Author struct {
	Name  string
	Email string
}
