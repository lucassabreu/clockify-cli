package cmdutil

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/spf13/cobra"
)

// AddProjectFlags creates a project flag with autocomplete configured
func AddProjectFlags(cmd *cobra.Command, f Factory) {
	cmd.Flags().StringP("project", "p", "",
		"the name/id of the project to work on")
	_ = cmd.MarkFlagRequired("project")

	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f))

}
