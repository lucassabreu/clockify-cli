package util

import (
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

	cmdutil.AddProjectFlags(cmd, f)
}

// FlagsDTO holds data about editing or creating a Task
type FlagsDTO struct {
	Workspace   string
	ProjectID   string
	Name        string
	Estimate    *time.Duration
	AssigneeIDs *[]string
	Billable    *bool
}

// TaskReadFlags read the common flags expected when editing a task
func TaskReadFlags(cmd *cobra.Command, f cmdutil.Factory) (p FlagsDTO, err error) {
	if err := cmdutil.XorFlag(map[string]bool{
		"assignee":    cmd.Flags().Changed("assignee"),
		"no-assignee": cmd.Flags().Changed("no-assignee"),
	}); err != nil {
		return p, err
	}

	if err := cmdutil.XorFlag(map[string]bool{
		"billable":     cmd.Flags().Changed("billable"),
		"not-billable": cmd.Flags().Changed("not-billable"),
	}); err != nil {
		return p, err
	}

	if p.Workspace, err = f.GetWorkspaceID(); err != nil {
		return
	}

	p.ProjectID, _ = cmd.Flags().GetString("project")
	p.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("estimate") {
		e, _ := cmd.Flags().GetInt32("estimate")
		d := time.Duration(e) * time.Hour
		p.Estimate = &d
	}

	if cmd.Flags().Changed("assignee") {
		assignees, _ := cmd.Flags().GetStringSlice("assignee")
		p.AssigneeIDs = &assignees
	}

	if f.Config().IsAllowNameForID() {
		c, err := f.Client()
		if err != nil {
			return p, err
		}

		if p.ProjectID, err = search.GetProjectByName(
			c, f.Config(), p.Workspace, p.ProjectID, ""); err != nil {
			return p, err
		}

		if p.AssigneeIDs != nil {
			as := *p.AssigneeIDs
			if as, err = search.GetUsersByName(
				c, p.Workspace, as); err != nil {
				return p, err
			}
			p.AssigneeIDs = &as
		}
	}

	if cmd.Flags().Changed("no-assignee") {
		var a []string

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
