package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/lucassabreu/clockify-cli/pkg/uiutil"
)

// GetDatesInteractiveFn will ask the user the start and end times of the entry
func GetDatesInteractiveFn(f cmdutil.Factory) Step {
	if !f.Config().IsInteractive() {
		return skip
	}

	return func(t TimeEntryDTO) (TimeEntryDTO, error) {
		return askTimeEntryDatesInteractive(f.UI(), t)
	}
}

func askTimeEntryDatesInteractive(
	ui ui.UI,
	dto TimeEntryDTO,
) (TimeEntryDTO, error) {
	var err error
	dateString := dto.Start.In(time.Local).
		Format(timehlp.FullTimeFormat)
	if dto.Start, err = ui.AskForDateTime(
		"Start", dateString, timehlp.ConvertToTime); err != nil {
		return dto, err
	}

	dateString = ""
	if dto.End != nil {
		dateString = dto.End.In(time.Local).
			Format(timehlp.FullTimeFormat)
	}

	if dto.End, err = ui.AskForDateTimeOrNil(
		"End", dateString, timehlp.ConvertToTime); err != nil {
		return dto, err
	}

	return dto, nil
}

// GetPropsInteractiveFn will return a callback that asks the user
// interactively about the properties of the time entry, only if the parameter
// cmdutil.CONF_INTERACTIVE is active
func GetPropsInteractiveFn(
	dc DescriptionSuggestFn,
	f cmdutil.Factory,
) Step {
	if !f.Config().IsInteractive() {
		return skip
	}

	return func(tei TimeEntryDTO) (TimeEntryDTO, error) {
		c, err := f.Client()
		if err != nil {
			return tei, err
		}

		return askTimeEntryPropsInteractive(
			tei,
			c,
			f.UI(),
			dc,
			f.Config().IsAllowArchivedTags(),
		)
	}
}

func askTimeEntryPropsInteractive(
	te TimeEntryDTO,
	c api.Client,
	ui ui.UI,
	dc DescriptionSuggestFn,
	allowArchived bool,
) (TimeEntryDTO, error) {
	var err error
	w, err := c.GetWorkspace(api.GetWorkspace{ID: te.Workspace})
	if err != nil {
		return te, err
	}

	te.ProjectID, err = getProjectID(te.ProjectID, w, c, ui)
	if err != nil {
		return te, err
	}

	if te.ProjectID != "" {
		te.TaskID, err = getTaskID(te.TaskID, te.ProjectID, w, c, ui)
		if err != nil {
			return te, err
		}
	}

	te.Description = getDescription(te.Description, dc, ui,
		w.Settings.ForceDescription)

	te.TagIDs, err = getTagIDs(te.TagIDs, w, c, allowArchived, ui)

	return te, err
}

func getProjectID(
	projectID string, w dto.Workspace, c api.Client, ui ui.UI,
) (string, error) {
	b := false
	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       w.ID,
		Archived:        &b,
		PaginationParam: api.AllPages(),
	})

	if err != nil || len(projects) == 0 {
		return "", err
	}

	p, err := uiutil.AskProject(uiutil.AskProjectParam{
		UI:            ui,
		ProjectID:     projectID,
		Projects:      projects,
		ForceProjects: w.Settings.ForceDescription,
	})

	if p != nil {
		return p.ID, err
	}

	return "", err
}

func getTaskID(
	taskID, projectID string, w dto.Workspace, c api.Client, ui ui.UI,
) (string, error) {
	tasks, err := c.GetTasks(api.GetTasksParam{
		Workspace:       w.ID,
		ProjectID:       projectID,
		PaginationParam: api.AllPages(),
		Active:          true,
	})

	// todo: this is a workaround for the cli, the api needs to be fixed
	var httpErr dto.Error
	if errors.As(err, &httpErr) && httpErr.Code == 501 && strings.Contains(
		httpErr.Message,
		"doesn't belong to PROJECT with id "+projectID,
	) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "", nil
	}

	t, err := uiutil.AskTask(uiutil.AskTaskParam{
		UI:     ui,
		TaskID: taskID,
		Tasks:  tasks,
		Force:  w.Settings.ForceTasks,
	})

	if t != nil {
		return t.ID, err
	}

	return "", err
}

func getDescription(
	description string,
	dc DescriptionSuggestFn,
	i ui.UI,
	force bool,
) string {
	var v func(string) error
	if force {
		v = func(s string) error {
			if s == "" {
				return errors.New("description should be informed")
			}
			return nil
		}
	}

	description, _ = i.AskForValidText(
		"Description:",
		v,
		ui.WithDefault(description),
		ui.WithSuggestion(dc),
	)
	return description
}

func getTagIDs(
	tagIDs []string, w dto.Workspace, c api.Client, allowArchived bool,
	ui ui.UI,
) ([]string, error) {
	var archived *bool
	if !allowArchived {
		f := false
		archived = &f
	}
	tags, err := c.GetTags(api.GetTagsParam{
		Workspace: w.ID,
		Archived:  archived,
	})

	if err != nil || len(tags) == 0 {
		return nil, err
	}

	tagsString := make([]string, len(tags))
	for i, u := range tags {
		tagsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
	}

	current := make([]string, len(tagIDs))
	for i, t := range tagIDs {
		for _, s := range tagsString {
			if strings.HasPrefix(s, t) {
				current[i] = s
				break
			}
		}
	}

	ts, err := uiutil.AskTags(uiutil.AskTagsParam{
		UI:     ui,
		TagIDs: tagIDs,
		Tags:   tags,
		Force:  w.Settings.ForceTags,
	})

	if err != nil || len(ts) == 0 {
		return nil, err
	}

	newTags := make([]string, len(ts))
	for i := range ts {
		newTags[i] = ts[i].ID
	}

	return newTags, err
}
