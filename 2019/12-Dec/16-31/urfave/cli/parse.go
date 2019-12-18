package cli

import (
	"flag"
	"strings"
)

type iterativeParser interface {
	newFlagSet() (*flag.FlagSet, error)
	useShortOptionHandling() bool
}

func parseIter(set *flag.FlagSet, ip iterativeParser, args []string, shellComplete bool) error {
	for {
		err := set.Parse(args)
		if !ip.useShortOptionHandling() || err == nil {
			if shellComplete {
				return nil
			}
			return err
		}

		errStr := err.Error()
		trimmed := strings.TrimPrefix(errStr, "flag provided but not defined: -")
		if errStr == trimmed {
			return err
		}

		argsWereSplit := false
		for i, arg := range args {
			if name := strings.TrimLeft(arg, "-"); name != trimmed {
				continue
			}

			shortOpts := splitShortOptions(set, arg)
			if len(shortOpts) == 1 {
				return err
			}

			args = append(args[:i], append(shortOpts, args[i+1:]...)...)
			argsWereSplit = true
			break
		}

		if !argsWereSplit {
			return err
		}

		newSet, err := ip.newFlagSet()
		if err != nil {
			return err
		}
		*set = *newSet
	}
}

func splitShortOptions(set *flag.FlagSet, arg string) []string {
	shortFlagsExist := func(s string) bool {
		for _, c := range s[1:] {
			if f := set.Lookup(string(c)); f == nil {
				return false
			}
		}
		return true
	}

	if !isSplittable(arg) || !shortFlagsExist(arg) {
		return []string{arg}
	}
	separated := make([]string, 0, len(arg)-1)
	for _, flagChar := range arg[1:] {
		separated = append(separated, "-"+string(flagChar))
	}
	return separated
}

func isSplittable(flagArg string) bool {
	return strings.HasPrefix(flagArg, "-") && !strings.HasPrefix(flagArg, "--") && len(flagArg) > 2
}
