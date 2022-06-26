package cmdutil

import (
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RequiredNamedArgs will fail if the number of arguments received is less than
// the length of the parameter `names`, and will show the help explaining
// required named arguments
func RequiredNamedArgs(names ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= len(names) {
			return nil
		}

		if len(names) == 1 {
			return FlagErrorWrap(errors.New("requires arg " + names[0]))
		}

		return FlagErrorWrap(errors.Errorf(
			"requires args %s; %d of those received",
			strhlp.ListForHumans(names), len(args),
		))
	}
}
