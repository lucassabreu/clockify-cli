package util

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
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
			f.Config().GetBool(cmdutil.CONF_ALLOW_ARCHIVED_TAGS),
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

const noProject = "No Project"

func getProjectID(
	projectID string, w dto.Workspace, c api.Client, ui ui.UI,
) (string, error) {
	b := false
	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       w.ID,
		Archived:        &b,
		PaginationParam: api.AllPages(),
	})

	if err != nil {
		return "", err
	}

	projectsString := make([]string, len(projects))
	found := -1
	projectNameSize := 0

	for i := range projects {
		projectsString[i] = projects[i].ID + " - " + projects[i].Name
		if c := utf8.RuneCountInString(projectsString[i]); projectNameSize < c {
			projectNameSize = c
		}

		if found == -1 && projects[i].ID == projectID {
			projectID = projectsString[i]
			found = i
		}
	}

	format := fmt.Sprintf("%%-%ds| %%s", projectNameSize+1)

	for i := range projects {
		client := "Without Client"
		if projects[i].ClientID != "" {
			client = "Client: " + projects[i].ClientName +
				" (" + projects[i].ClientID + ")"
		}

		projectsString[i] = fmt.Sprintf(
			format,
			projectsString[i],
			client,
		)
	}

	if found == -1 {
		if projectID != "" {
			fmt.Printf("Project '%s' informed was not found.\n", projectID)
			projectID = ""
		}
	} else {
		projectID = projectsString[found]
	}

	if !w.Settings.ForceProjects {
		projectsString = append([]string{noProject}, projectsString...)
	}

	projectID, err = ui.AskFromOptions("Choose your project:",
		projectsString, projectID)
	if err != nil || projectID == noProject || projectID == "" {
		return "", err
	}

	return strings.TrimSpace(projectID[0:strings.Index(projectID, " - ")]), nil
}

const noTask = "No Task"

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

	tasksString := make([]string, len(tasks))
	found := -1

	for i := range tasks {
		tasksString[i] = tasks[i].ID + " - " + tasks[i].Name

		if found == -1 && tasks[i].ID == taskID {
			taskID = tasksString[i]
			found = i
		}
	}

	if found == -1 {
		if taskID != "" {
			fmt.Printf("Task '%s' informed was not found.\n", taskID)
			taskID = ""
		}
	} else {
		taskID = tasksString[found]
	}

	if !w.Settings.ForceTasks {
		tasksString = append([]string{noTask}, tasksString...)
	}

	taskID, err = ui.AskFromOptions("Choose your task:", tasksString, taskID)
	if err != nil || taskID == noTask || taskID == "" {
		return "", err
	}

	return strings.TrimSpace(taskID[0:strings.Index(taskID, " - ")]), nil
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
				return errors.New("description should be informed !")
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

	if err != nil {
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

	var newTags []string
	if newTags, err = ui.AskManyFromOptions("Choose your tags:",
		tagsString, current, func(s []string) error {
			if w.Settings.ForceTags && len(s) == 0 {
				return errors.New("at least one tag should be selected")
			}

			return nil
		}); err != nil {
		return nil, nil
	}

	for i, t := range newTags {
		newTags[i] = strings.TrimSpace(t[0:strings.Index(t, " - ")])
	}

	return newTags, nil
}
