package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"
)

var fullTimeFormat = "2006-01-02 15:04:05"
var simplerTimeFormat = "2006-01-02 15:04"
var onlyTimeFormat = "15:04:05"
var simplerOnlyTimeFormat = "15:04"
var nowTimeFormat = "now"

func withClockifyClient(fn func(cmd *cobra.Command, args []string, c *api.Client) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := getAPIClient()
		if err != nil {
			return err
		}

		return fn(cmd, args, c)
	}
}

func convertToTime(timeString string) (t time.Time, err error) {
	timeString = strings.TrimSpace(timeString)

	if nowTimeFormat == strings.ToLower(timeString) {
		return time.Now().In(time.Local), nil
	}

	if len(fullTimeFormat) != len(timeString) && len(simplerTimeFormat) != len(timeString) && len(onlyTimeFormat) != len(timeString) && len(simplerOnlyTimeFormat) != len(timeString) {
		return t, fmt.Errorf(
			"supported formats are: %s",
			strings.Join(
				[]string{fullTimeFormat, simplerTimeFormat, onlyTimeFormat, simplerOnlyTimeFormat, nowTimeFormat},
				", ",
			),
		)
	}

	if len(simplerOnlyTimeFormat) == len(timeString) || len(simplerTimeFormat) == len(timeString) {
		timeString = timeString + ":00"
	}

	if len(onlyTimeFormat) == len(timeString) {
		timeString = time.Now().Format("2006-01-02") + " " + timeString
	}

	return time.ParseInLocation(fullTimeFormat, timeString, time.Local)
}

func getAPIClient() (*api.Client, error) {
	c, err := api.NewClient(viper.GetString("token"))
	if err != nil {
		return c, err
	}

	if viper.GetBool("debug") {
		c.SetDebugLogger(
			log.New(os.Stdout, "DEBUG ", log.LstdFlags),
		)
	}

	return c, err
}

func getDateTimeParam(name string, required bool, value string, convert func(string) (time.Time, error)) (*time.Time, error) {
	var t time.Time
	var err error

	message := fmt.Sprintf("%s (leave it blank for empty):", name)
	if required {
		message = fmt.Sprintf("%s:", name)
	}

	for {
		_ = survey.AskOne(
			&survey.Input{
				Message: message,
				Default: value,
			},
			&value,
			nil,
		)

		if value == "" && !required {
			return nil, nil
		}

		if t, err = convertToTime(value); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		return &t, err
	}
}

func newEntry(c *api.Client, te dto.TimeEntryImpl, interactive, autoClose bool, format string, asJSON bool) error {
	var err error

	if interactive {
		te.ProjectID, err = getProjectID(te.ProjectID, te.WorkspaceID, c)
		if err != nil {
			return err
		}
	}

	if te.ProjectID == "" {
		return errors.New("project must be informed")
	}

	if interactive {
		te.Description = getDescription(te.Description)

		te.TagIDs, err = getTagIDs(te.TagIDs, te.WorkspaceID, c)
		if err != nil {
			return err
		}

		var date *time.Time
		dateString := te.TimeInterval.Start.Format(fullTimeFormat)

		if date, err = getDateTimeParam("Start", true, dateString, convertToTime); err != nil {
			return err
		}
		te.TimeInterval.Start = *date

		dateString = ""
		if te.TimeInterval.End != nil {
			dateString = te.TimeInterval.End.Format(fullTimeFormat)
		}

		if date, err = getDateTimeParam("End", false, dateString, convertToTime); err != nil {
			return err
		}
		te.TimeInterval.End = date
	}

	if autoClose {
		err = c.Out(api.OutParam{
			Workspace: te.WorkspaceID,
			End:       te.TimeInterval.Start,
		})

		if err != nil {
			return err
		}
	}

	tei, err := c.CreateTimeEntry(api.CreateTimeEntryParam{
		Workspace:   te.WorkspaceID,
		Billable:    te.Billable,
		Start:       te.TimeInterval.Start,
		End:         te.TimeInterval.End,
		ProjectID:   te.ProjectID,
		Description: te.Description,
		TagIDs:      te.TagIDs,
		TaskID:      te.TaskID,
	})

	if err != nil {
		return err
	}

	fte, err := c.ConvertIntoFullTimeEntry(tei)
	if err != nil {
		return err
	}

	var reportFn func(*dto.TimeEntry, io.Writer) error

	reportFn = reports.TimeEntryPrint

	if asJSON {
		reportFn = reports.TimeEntryJSONPrint
	}

	if format != "" {
		reportFn = reports.TimeEntryPrintWithTemplate(format)
	}

	return reportFn(&fte, os.Stdout)
}

func getProjectID(projectID string, workspace string, c *api.Client) (string, error) {
	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace: workspace,
	})

	if err != nil {
		return "", err
	}

	projectsString := make([]string, len(projects))
	found := false
	for i, u := range projects {
		projectsString[i] = fmt.Sprintf("%s - %s", u.ID, u.Name)
		if u.ID == projectID {
			projectID = projectsString[i]
			found = true
		}
	}

	if !found {
		fmt.Printf("Project '%s' informed was not found.\n", projectID)
		projectID = ""
	}

	err = survey.AskOne(
		&survey.Select{
			Message: "Choose your project:",
			Options: projectsString,
			Default: projectID,
		},
		&projectID,
		nil,
	)

	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(projectID[0:strings.Index(projectID, " - ")]), nil
}

func getDescription(description string) string {
	_ = survey.AskOne(
		&survey.Input{
			Message: "Description:",
			Default: description,
		},
		&description,
		nil,
	)

	return description
}

func getTagIDs(tagIDs []string, workspace string, c *api.Client) ([]string, error) {
	if len(tagIDs) > 0 && !viper.GetBool("interactive") {
		return tagIDs, nil
	}

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

	err = survey.AskOne(
		&survey.MultiSelect{
			Message: "Choose your tags:",
			Options: tagsString,
			Default: tagIDs,
		},
		&tagIDs,
		nil,
	)

	if err != nil {
		return nil, nil
	}

	for i, t := range tagIDs {
		tagIDs[i] = strings.TrimSpace(t[0:strings.Index(t, " - ")])
	}

	return tagIDs, nil
}

func getUserId(c *api.Client) (string, error) {
	userId := viper.GetString("user.id")
	if len(userId) > 0 {
		return userId, nil
	}

	u, err := c.GetMe()
	if err != nil {
		return "", err
	}

	return u.ID, nil
}
