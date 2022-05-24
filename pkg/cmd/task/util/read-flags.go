package util

import (
	"errors"
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/spf13/cobra"
)

// TaskAddPropFlags add common flags expected on editing a task
func TaskAddPropFlags(cmd *cobra.Command, f cmdutil.Factory) {
	cmd.Flags().StringP("name", "n", "", "new name of the task")
	cmd.Flags().Int32P("estimate", "E", 0, "estimation on hours")
	cmd.Flags().Bool("billable", false, "sets the task as billable")
	cmd.Flags().Bool("not-billable", false, "sets the task as not billable")

	cmd.Flags().StringSliceP("assignee", "A", []string{},
		"list of users that are assigned to this task")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "assignee",
		cmdcomplutil.NewUserAutoComplete(f))

	cmd.Flags().Bool("no-assignee", false,
		"cleans the assignee list")
}

type FlagsDTO struct {
	Workspace   string
	Name        string
	Estimate    *time.Duration
	AssigneeIDs *[]string
	Billable    *bool
}

// TaskReadFlags read the common flags expected when editing a task
func TaskReadFlags(cmd *cobra.Command, f cmdutil.Factory) (p FlagsDTO, err error) {
	if p.Workspace, err = f.GetWorkspaceID(); err != nil {
		return
	}

	p.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("estimate") {
		e, _ := cmd.Flags().GetInt32("estimate")
		d := time.Duration(e) * time.Hour
		p.Estimate = &d
	}

	if cmd.Flags().Changed("assignee") && cmd.Flags().Changed("no-assignee") {
		return p, errors.New(
			"`--assignee` and `--no-assignee` can't be used together")
	}

	if cmd.Flags().Changed("assignee") {
		c, err := f.Client()
		if err != nil {
			return p, err
		}

		assignees, _ := cmd.Flags().GetStringSlice("assignee")
		if assignees, err = search.GetUsersByName(
			c, p.Workspace, assignees); err != nil {
			return p, err
		}

		p.AssigneeIDs = &assignees
	}

	if cmd.Flags().Changed("no-assignee") {
		a := []string{}
		p.AssigneeIDs = &a
	}

	switch {
	case cmd.Flags().Changed("billable"):
		b := true
		p.Billable = &b
	case cmd.Flags().Changed("not-billable"):
		b := false
		p.Billable = &b
	}

	return
}
