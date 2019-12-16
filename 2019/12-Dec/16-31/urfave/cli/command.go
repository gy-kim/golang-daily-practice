package cli

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string

	// A list of aliases for the command
	Aliases []string

	Usage                  string
	Description            string
	ArgsUsage              string
	Category               string
	BashComplete           BashCompleteFunc
	Before                 BeforeFunc
	After                  AfterFunc
	Action                 ActionFunc
	OnUsageError           OnUsageErrorFunc
	Subcommands            []*Command
	Flags                  []Flag
	SkipFlagParsing        bool
	HideHelp               bool
	Hidden                 bool
	UseShortOptionHandling bool
	HelpName               string
	commandNamePath        []string
	CustomHelpTemplate     string
}
