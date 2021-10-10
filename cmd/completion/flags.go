package completion

import (
	"github.com/spf13/cobra"
)

// AddFixedSuggestionsToFlag add fixed suggestions to a flag
func AddFixedSuggestionsToFlag(cmd *cobra.Command, flagName string, va ValidArgs) error {
	f := cmd.Flags().Lookup(flagName)
	f.Usage = va.IntoUse() + " " + f.Usage
	return cmd.RegisterFlagCompletionFunc(
		f.Name,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return va.IntoValidArgs(), cobra.ShellCompDirectiveDefault
		},
	)
}

type suggestFn func(cmd *cobra.Command, args []string, toComplete string) (ValidArgs, error)

func process(va ValidArgs, err error) ([]string, cobra.ShellCompDirective) {
	if err != nil {
		cobra.CompError(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}

	return va.IntoValidArgs(), cobra.ShellCompDirectiveDefault
}

// AddSuggestionsToFlag add fixed suggestions to a flag
func AddSuggestionsToFlag(cmd *cobra.Command, flagName string, suggestFn suggestFn) error {
	return cmd.RegisterFlagCompletionFunc(
		flagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return process(suggestFn(cmd, args, toComplete))
		},
	)
}

func EmptySuggestionFuncion(_ *cobra.Command, _ []string, _ string) (ValidArgs, error) {
	return EmptyValidArgs(), nil
}

// CombineSuggestionsToArgs combine one or more suggestion resolver functions and call then accordingly with arg count
func CombineSuggestionsToArgs(fns ...suggestFn) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > len(fns) {
			return []string{}, cobra.ShellCompDirectiveDefault
		}

		return process(fns[len(args)](cmd, args, toComplete))
	}
}
