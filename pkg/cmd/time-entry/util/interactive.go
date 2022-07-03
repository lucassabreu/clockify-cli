package util

import (
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
func GetDatesInteractiveFn(config cmdutil.Config) DoFn {
	if config.IsInteractive() {
		return askTimeEntryDatesInteractive
	}

	return nullCallback
}

func askTimeEntryDatesInteractive(
	te dto.TimeEntryImpl,
) (dto.TimeEntryImpl, error) {
	var err error
	dateString := te.TimeInterval.Start.In(time.Local).
		Format(timehlp.FullTimeFormat)
	if te.TimeInterval.Start, err = ui.AskForDateTime(
		"Start", dateString, timehlp.ConvertToTime); err != nil {
		return te, err
	}

	dateString = ""
	if te.TimeInterval.End != nil {
		dateString = te.TimeInterval.End.In(time.Local).
			Format(timehlp.FullTimeFormat)
	}

	if te.TimeInterval.End, err = ui.AskForDateTimeOrNil(
		"End", dateString, timehlp.ConvertToTime); err != nil {
		return te, err
	}

	return te, nil
}

// GetPropsInteractiveFn will return a callback that asks the user
// interactively about the properties of the time entry, only if the parameter
// cmdutil.CONF_INTERACTIVE is active
func GetPropsInteractiveFn(
	c *api.Client,
	dc DescriptionSuggestFn,
	config cmdutil.Config,
) DoFn {
	if config.IsInteractive() {
		return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			return askTimeEntryPropsInteractive(c, tei, dc)
		}
	}

	return nullCallback
}

func askTimeEntryPropsInteractive(
	c *api.Client,
	te dto.TimeEntryImpl,
	dc DescriptionSuggestFn,
) (dto.TimeEntryImpl, error) {
	var err error
	w, err := c.GetWorkspace(api.GetWorkspace{ID: te.WorkspaceID})
	if err != nil {
		return te, err
	}

	te.ProjectID, err = getProjectID(te.ProjectID, w, c)
	if err != nil {
		return te, err
	}

	if te.ProjectID != "" {
		te.TaskID, err = getTaskID(te.TaskID, te.ProjectID, w, c)
		if err != nil {
			return te, err
		}
	}

	te.Description = getDescription(te.Description, dc)

	te.TagIDs, err = getTagIDs(te.TagIDs, te.WorkspaceID, c)

	return te, err
}

const noProject = "No Project"

func getProjectID(
	projectID string, w dto.Workspace, c *api.Client) (string, error) {
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

	for i, u := range projects {
		projectsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
		if c := utf8.RuneCountInString(projectsString[i]); projectNameSize < c {
			projectNameSize = c
		}

		if found == -1 && u.ID == projectID {
			projectID = projectsString[i]
			found = i
		}
	}

	format := fmt.Sprintf("%%-%ds| %%s", projectNameSize+1)

	for i, u := range projects {
		client := "Without Client"
		if u.ClientID != "" {
			client = fmt.Sprintf("Client: %s (%s)", u.ClientName, u.ClientID)
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
	taskID, projectID string, w dto.Workspace, c *api.Client) (string, error) {
	tasks, err := c.GetTasks(api.GetTasksParam{
		Workspace:       w.ID,
		ProjectID:       projectID,
		PaginationParam: api.AllPages(),
		Active:          true,
	})

	if err != nil {
		return "", err
	}

	if len(tasks) == 0 {
		return "", nil
	}

	tasksString := make([]string, len(tasks))
	found := -1

	for i, u := range tasks {
		tasksString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)

		if found == -1 && u.ID == taskID {
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

func getDescription(description string, dc DescriptionSuggestFn) string {
	description, _ = ui.AskForText("Description:",
		ui.WithDefault(description),
		ui.WithSuggestion(dc))
	return description
}

func getTagIDs(
	tagIDs []string, workspace string, c *api.Client) ([]string, error) {
	tags, err := c.GetTags(api.GetTagsParam{
		Workspace: workspace,
	})

	if err != nil {
		return nil, err
	}

	tagsString := make([]string, len(tags))
	for i, u := range tags {
		tagsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
	}

	for i, t := range tagIDs {
		for _, s := range tagsString {
			if strings.HasPrefix(s, t) {
				tagIDs[i] = s
				break
			}
		}
	}

	var newTags []string
	if newTags, err = ui.AskManyFromOptions("Choose your tags:",
		tagsString, tagIDs); err != nil {
		return nil, nil
	}

	for i, t := range newTags {
		newTags[i] = strings.TrimSpace(t[0:strings.Index(t, " - ")])
	}

	return newTags, nil
}
