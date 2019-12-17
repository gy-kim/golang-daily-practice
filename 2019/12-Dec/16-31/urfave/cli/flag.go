package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const defaultPlaceholder = "value"

var (
	slPfx = fmt.Sprintf("sl:::%d:::", time.Now().UTC().UnixNano())

	commaWhitespace = regexp.MustCompile("[, ]+.*")
)

// BashCompletionFlag enables bash-completion for all commands and subcommands
var BashCompletionFlag Flag = &BoolFlag{
	Name:   "generate-bash-completion",
	Hidden: true,
}

// VersionFlag prints the version for the application
var VersionFlag Flag = &BoolFlag{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "print the version",
}

// HelpFlag prints the help for all commands and subcommands.
var HelpFlag Flag = &BoolFlag{
	Name:    "help",
	Aliases: []string{"h"},
	Usage:   "show help",
}

// FlagStringer converts a flag definition to a string. This is used by help to display a flag.
var FlagStringer FlagStringFunc = stringifyFlag

// Serializer is used to circumvent the limitations of flag.FlagSet.Set
type Serializer interface {
	Serializer() string
}

// FlagNamePrefixer converts a full flag name and its placeholer into the help
// message falg prefix. This is used by the default FlagStringer.
var FlagNamePrefixer FlagNamePrefixFunc = prefixedNames

// FlagEnvHinter annotates flag help message with the environment variable details.
// This is used by the default FlagStringer.
var FlagEnvHinter FlagEnvHintFunc = withEnvHint

// FlagFileHinter annotates flag help message with the environment variable details.
// This is used by the default FlagStringer.
var FlagFileHinter FlagFileHintFunc = withFileHint

// FlagsByName is a slice of Flag.
type FlagsByName []Flag

func (f FlagsByName) Len() int {
	return len(f)
}

func (f FlagsByName) Less(i, j int) bool {
	if len(f[j].Names()) == 0 {
		return false
	} else if len(f[i].Names()) == 0 {
		return true
	}

	return lexicographicLess(f[i].Names()[0], f[j].Names()[0])
}

func (f FlagsByName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Flag is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recommended that
// this interface be implemented.
type Flag interface {
	fmt.Stringer
	// Apply Flag settings to the given flag set
	Apply(*flag.FlagSet) error
	Names() []string
	IsSet() bool
}

// RequiredFlag is an interface that allows us to mark flags as required
// it allows flags required flags to be backwards compatible with the Flag interface
type RequiredFlag interface {
	Flag

	IsRequired() bool
}

// DocGenerationFlag is an interface that allows documentation generation for the flag
type DocGenerationFlag interface {
	Flag

	// TakesValue returns true if the falg takes a value, otherwise false
	TakesValue() bool

	// GetUsage returns the usage string for the flag
	GetUsage() string

	// GetValue returns the flags value as string representation and an empty
	// string if the flag takes no value at all.
	GetValue() string
}

func flagSet(name string, flags []Flag) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		if err := f.Apply(set); err != nil {
			return nil, err
		}
	}

	set.SetOutput(ioutil.Discard)
	return set, nil
}

func visibleFlags(fl []Flag) []Flag {
	var visible []Flag
	for _, f := range fl {
		field := flagValue(f).FieldByName("Hidden")
		if !field.IsValid() || !field.Bool() {
			visible = append(visible, f)
		}
	}
	return visible
}

func stringifyStringSliceFlag(f *StringSliceFlag) string {
	var defaultVals []string
	if f.Value != nil && len(f.Value.Value()) > 0 {
		for _, s := range f.Value.Value() {
			if len(s) > 0 {
				defaultVals = append(defaultVals, strconv.Quote(s))
			}
		}
	}
	return stringifySliceFlag(f.Usage, f.Names(), defaultVals)
}

func stringifySliceFlag(usage string, names, defaultVals []string) string {
	placeholder, usage := unquoteUage(usage)
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}

	defaultVal := ""
	if len(defaultVals) > 0 {
		defaultVal = fmt.Sprintf(" (default: %s)", strings.Join(defaultVals, ", "))
	}

	usageWithDefault := strings.TrimSpace(fmt.Sprintf("%s%s", usage, defaultVal))
	return fmt.Sprintf("%s\n%s", prefixedNames(names, placeholder), usageWithDefault)
}

func hasFlag(flags []Flag, fl Flag) bool {
	for _, existing := range flags {
		if fl == existing {
			return true
		}
	}
	return false
}

func flagFromEnvOrFile(envVars []string, filePath string) (val string, ok bool) {
	for _, envVar := range envVars {
		envVar = strings.TrimSpace(envVar)
		if val, ok := syscall.Getenv(envVar); ok {
			return val, true
		}
	}
	for _, fileVar := range strings.Split(filePath, ",") {
		if data, err := ioutil.ReadFile(fileVar); err == nil {
			return string(data), true
		}
	}
	return "", false
}
