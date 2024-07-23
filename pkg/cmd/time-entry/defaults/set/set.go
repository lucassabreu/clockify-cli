package set

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	outd "github.com/lucassabreu/clockify-cli/pkg/output/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/lucassabreu/clockify-cli/pkg/uiutil"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdSet sets the default parameters for time entries in the current folder
func NewCmdSet(
	f cmdutil.Factory,
	report func(outd.OutputFlags, io.Writer, defaults.DefaultTimeEntry) error,
) *cobra.Command {
	if report == nil {
		panic(errors.New("report parameter should not be nil"))
	}

	short := "Sets the default parameters for the current folder"
	of := outd.OutputFlags{}
	cmd := &cobra.Command{
		Use:   "set",
		Short: short,
		Long: short + "\n" +
			"The parameters will be saved in the current working directory " +
			"in the file " + defaults.DEFAULT_FILENAME + ".yaml",
		Example: "",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.XorFlagSet(
				cmd.Flags(), "billable", "not-billable"); err != nil {
				return err
			}

			d, err := f.TimeEntryDefaults().Read()
			if err != nil && err != defaults.DefaultsFileNotFoundErr {
				return err
			}

			n, changed := readFlags(d, cmd.Flags())

			var w string
			if w, err = f.GetWorkspaceID(); err != nil {
				return err
			}

			c, err := f.Client()
			if err != nil {
				return err
			}

			if changed {
				if n.TaskID != "" && n.ProjectID == "" {
					return errors.New("can't set task without project")
				}

				if f.Config().IsAllowNameForID() {
					if n, err = updateIDsByNames(
						c, n, f.Config(), w); err != nil {
						return err
					}
				} else {
					if err = checkIDs(c, w, n); err != nil {
						return err
					}
				}
			}

			if f.Config().IsInteractive() {
				if n, err = ask(n, w, f.Config(), c, f.UI()); err != nil {
					return err
				}
			}

			if err = f.TimeEntryDefaults().Write(n); err != nil {
				return err
			}

			return report(of, cmd.OutOrStdout(), n)
		},
	}

	cmd.Flags().StringVarP(&of.Format,
		"format", "f", outd.FORMAT_YAML, "output format")
	_ = cmdcompl.AddFixedSuggestionsToFlag(cmd, "format",
		cmdcompl.ValidArgsSlide{outd.FORMAT_YAML, outd.FORMAT_JSON})

	cmd.Flags().BoolP("billable", "b", false,
		"time entry should be billable by default")
	cmd.Flags().BoolP("not-billable", "n", false,
		"time entry should not be billable by default")
	cmd.Flags().String("task", "", "default task")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "task",
		cmdcomplutil.NewTaskAutoComplete(f, true))

	cmd.Flags().StringSliceP("tag", "T", []string{},
		"add tags be used by default")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "tag",
		cmdcomplutil.NewTagAutoComplete(f))

	cmd.Flags().StringP("project", "p", "", "project to used by default")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f, f.Config()))

	return cmd
}

func readFlags(
	d defaults.DefaultTimeEntry,
	f *pflag.FlagSet,
) (defaults.DefaultTimeEntry, bool) {
	changed := false
	if f.Changed("project") {
		d.ProjectID, _ = f.GetString("project")
		changed = true
	}

	if f.Changed("task") {
		d.TaskID, _ = f.GetString("task")
		changed = true
	}

	if f.Changed("tag") {
		d.TagIDs, _ = f.GetStringSlice("tag")
		d.TagIDs = strhlp.Unique(d.TagIDs)
		changed = true
	}

	if f.Changed("billable") {
		b := true
		d.Billable = &b
		changed = true
	} else if f.Changed("not-billable") {
		b := false
		d.Billable = &b
		changed = true
	}

	return d, changed
}

func checkIDs(c api.Client, w string, d defaults.DefaultTimeEntry) error {
	if d.ProjectID != "" {
		p, err := c.GetProject(api.GetProjectParam{
			Workspace: w,
			ProjectID: d.ProjectID,
			Hydrate:   d.TaskID != "",
		})

		if err != nil {
			return err
		}

		if d.TaskID != "" {
			found := false
			for i := range p.Tasks {
				if p.Tasks[i].ID == d.TaskID {
					found = true
					break
				}
			}

			if !found {
				return errors.New(
					"can't find task with ID \"" + d.TaskID +
						"\" on project \"" + d.ProjectID + "\"")
			}
		}
	} else if d.TaskID != "" {
		return errors.New("task can't be set without a project")
	}

	tags, err := c.GetTags(api.GetTagsParam{
		Workspace:       w,
		Archived:        &archived,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return err
	}

	ids := make([]string, len(tags))
	for i := range tags {
		ids[i] = tags[i].ID
	}

	for _, id := range d.TagIDs {
		if !strhlp.InSlice(id, ids) {
			return errors.Errorf("can't find tag with ID \"%s\"", id)
		}
	}

	return nil
}

var archived = false

func updateIDsByNames(
	c api.Client, d defaults.DefaultTimeEntry,
	cnf cmdutil.Config, w string) (
	defaults.DefaultTimeEntry,
	error,
) {
	var err error
	if d.ProjectID != "" {
		d.ProjectID, err = search.GetProjectByName(c, cnf,
			w, d.ProjectID, "")
		if err != nil {
			d.ProjectID = ""
			d.TaskID = ""
			if !cnf.IsInteractive() {
				return d, err
			}
		}
	}

	if d.TaskID != "" {
		d.TaskID, err = search.GetTaskByName(c, api.GetTasksParam{
			Workspace: w,
			ProjectID: d.ProjectID,
			Active:    true,
		}, d.TaskID)
		if err != nil && !cnf.IsInteractive() {
			return d, err
		}
	}

	if len(d.TagIDs) > 0 {
		d.TagIDs, err = search.GetTagsByName(
			c, w, !cnf.IsAllowArchivedTags(), d.TagIDs)
		if err != nil && !cnf.IsInteractive() {
			return d, err
		}
	}

	return d, nil
}

func ask(
	d defaults.DefaultTimeEntry,
	w string,
	cnf cmdutil.Config,
	c api.Client,
	ui ui.UI,
) (
	defaults.DefaultTimeEntry,
	error,
) {
	ui.SetPageSize(uint(cnf.InteractivePageSize()))

	ps, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       w,
		Archived:        &archived,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return d, err
	}

	p, err := uiutil.AskProject(uiutil.AskProjectParam{
		UI:        ui,
		ProjectID: d.ProjectID,
		Projects:  ps,
	})
	if err != nil {
		return d, err
	}
	if p != nil {
		d.ProjectID = p.ID
	} else {
		d.ProjectID = ""
	}

	if d.ProjectID != "" {
		ts, err := c.GetTasks(api.GetTasksParam{
			Workspace:       w,
			ProjectID:       d.ProjectID,
			Active:          true,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return d, err
		}

		t, err := uiutil.AskTask(uiutil.AskTaskParam{
			UI:     ui,
			TaskID: d.TaskID,
			Tasks:  ts,
		})
		if err != nil {
			return d, err
		}
		if t != nil {
			d.TaskID = t.ID
		} else {
			d.TaskID = ""
		}
	} else {
		d.TaskID = ""
	}

	var archived *bool
	if !cnf.IsAllowArchivedTags() {
		b := false
		archived = &b
	}

	tags, err := c.GetTags(api.GetTagsParam{
		Workspace:       w,
		Archived:        archived,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return d, err
	}

	tags, err = uiutil.AskTags(uiutil.AskTagsParam{
		UI:     ui,
		TagIDs: d.TagIDs,
		Tags:   tags,
	})
	if err != nil {
		return d, err
	}
	d.TagIDs = make([]string, len(tags))
	for i := range tags {
		d.TagIDs[i] = tags[i].ID
	}

	b := false
	if d.Billable != nil {
		b = *d.Billable
	}

	b, err = ui.Confirm("Should be billable?", b)
	if err != nil {
		return d, err
	}
	d.Billable = &b

	return d, err
}
