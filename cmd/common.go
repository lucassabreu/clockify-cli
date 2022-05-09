package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/lucassabreu/clockify-cli/internal/output"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/lucassabreu/clockify-cli/ui"
	stackedErrors "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	c, err := api.NewClient(viper.GetString(TOKEN))
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

func getTagsByNameOrId(c *api.Client, workspace string, tags []string) ([]string, error) {
	dtos, err := c.GetTags(api.GetTagsParam{
		Workspace:       workspace,
		PaginationParam: api.PaginationParam{AllPages: true},
	})

	if err != nil {
		return tags, err
	}

	for i, id := range tags {
		id = strhlp.Normalize(strings.TrimSpace(id))
		found := false
		for _, dto := range dtos {
			if strings.ToLower(dto.ID) == id {
				tags[i] = dto.ID
				found = true
				break
			}

			if strings.Contains(strhlp.Normalize(dto.Name), id) {
				tags[i] = dto.ID
				found = true
				break
			}
		}

		if !found {
			return tags, stackedErrors.Errorf("No tag with id or name containing: %s", id)
		}
	}

	return tags, nil
}

func getClientByNameOrId(c *api.Client, workspace string, client string) (string, error) {
	client = strhlp.Normalize(strings.TrimSpace(client))
	clients, err := c.GetClients(api.GetClientsParam{
		Workspace:       workspace,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return "", err
	}

	for _, cl := range clients {
		if strings.ToLower(cl.ID) == client {
			return cl.ID, nil
		}
		if strings.Contains(strhlp.Normalize(cl.Name), client) {
			return cl.ID, nil
		}
	}

	return "", stackedErrors.Errorf("No client with id or name containing: %s", client)
}

func getProjectByNameOrId(c *api.Client, workspace, project string) (string, error) {
	project = strhlp.Normalize(strings.TrimSpace(project))
	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       workspace,
		PaginationParam: api.PaginationParam{AllPages: true},
	})
	if err != nil {
		return "", err
	}

	for _, p := range projects {
		if strings.ToLower(p.ID) == project {
			return p.ID, nil
		}
		if strings.Contains(strhlp.Normalize(p.Name), project) {
			return p.ID, nil
		}
	}

	return "", stackedErrors.Errorf("No project with id or name containing: %s", project)
}

func getTaskByNameOrId(c *api.Client, workspace, project, task string) (string, error) {
	task = strhlp.Normalize(strings.TrimSpace(task))
	tasks, err := c.GetTasks(api.GetTasksParam{
		Workspace:       workspace,
		ProjectID:       project,
		PaginationParam: api.PaginationParam{AllPages: true},
	})
	if err != nil {
		return "", err
	}

	for _, p := range tasks {
		if strings.ToLower(p.ID) == task {
			return p.ID, nil
		}
		if strings.Contains(strhlp.Normalize(p.Name), task) {
			return p.ID, nil
		}
	}

	return "", stackedErrors.Errorf("No task with id or name containing: %s", task)
}

func askTimeEntryPropsInteractive(
	c *api.Client,
	te dto.TimeEntryImpl,
	dc suggestFn,
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

func askTimeEntryDatesInteractive(
	te dto.TimeEntryImpl,
) (dto.TimeEntryImpl, error) {
	var err error
	dateString := te.TimeInterval.Start.In(time.Local).Format(fullTimeFormat)
	if te.TimeInterval.Start, err = ui.AskForDateTime("Start", dateString, convertToTime); err != nil {
		return te, err
	}

	dateString = ""
	if te.TimeInterval.End != nil {
		dateString = te.TimeInterval.End.In(time.Local).Format(fullTimeFormat)
	}

	if te.TimeInterval.End, err = ui.AskForDateTimeOrNil("End", dateString, convertToTime); err != nil {
		return te, err
	}

	return te, nil
}

func validateTimeEntry(te dto.TimeEntryImpl, c *api.Client) error {
	w, err := c.GetWorkspace(api.GetWorkspace{ID: te.WorkspaceID})
	if err != nil {
		return err
	}

	if w.Settings.ForceProjects && te.ProjectID == "" {
		return errors.New("workspace requires project")
	}

	if w.Settings.ForceDescription && strings.TrimSpace(te.Description) == "" {
		return errors.New("workspace requires description")
	}

	if w.Settings.ForceTags && len(te.TagIDs) == 0 {
		return errors.New("workspace requires at least one tag")
	}

	if te.ProjectID == "" {
		return nil
	}

	p, err := c.GetProject(api.GetProjectParam{
		Workspace: te.WorkspaceID,
		ProjectID: te.ProjectID,
	})

	if err != nil {
		return err
	}

	if p.Archived {
		return fmt.Errorf("project %s - %s is archived", p.ID, p.Name)
	}

	return nil
}

func printTimeEntryImpl(
	tei dto.TimeEntryImpl,
	c *api.Client,
	cmd *cobra.Command,
	timeFormat string,
) error {
	fte, err := c.GetHydratedTimeEntry(api.GetTimeEntryParam{
		Workspace:   tei.WorkspaceID,
		TimeEntryID: tei.ID,
	})
	if err != nil {
		return err
	}

	return printTimeEntry(fte, cmd, timeFormat)
}

type CallbackFn func(dto.TimeEntryImpl) (dto.TimeEntryImpl, error)

func nullCallback(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
	return te, nil
}

func lookupProject(c *api.Client) CallbackFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if te.ProjectID == "" {
			return te, nil
		}

		var err error
		te.ProjectID, err = getProjectByNameOrId(c, te.WorkspaceID, te.ProjectID)
		return te, err
	}

}

func lookupTask(c *api.Client) CallbackFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if te.TaskID == "" {
			return te, nil
		}

		var err error
		te.TaskID, err = getTaskByNameOrId(c,
			te.WorkspaceID, te.ProjectID, te.TaskID)
		return te, err
	}
}

func lookupTags(c *api.Client) CallbackFn {
	return func(te dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		if len(te.TagIDs) == 0 {
			return te, nil
		}

		var err error
		te.TagIDs, err = getTagsByNameOrId(c, te.WorkspaceID, te.TagIDs)
		return te, err
	}

}

func disableErrorReporting(cbs []CallbackFn) []CallbackFn {
	for i := range cbs {
		cb := cbs[i]
		cbs[i] = func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			tei, _ = cb(tei)
			return tei, nil
		}
	}
	return cbs
}

func composeCallbacks(cbs ...CallbackFn) CallbackFn {
	return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		var err error
		for _, cb := range cbs {
			if tei, err = cb(tei); err != nil {
				return tei, err
			}
		}

		return tei, err
	}
}

func getAllowNameForIDsFn(c *api.Client) CallbackFn {
	if !viper.GetBool(ALLOW_NAME_FOR_ID) {
		return nullCallback
	}

	cbs := []CallbackFn{
		lookupProject(c),
		lookupTask(c),
		lookupTags(c),
	}

	if viper.GetBool(INTERACTIVE) {
		cbs = disableErrorReporting(cbs)
	}

	return composeCallbacks(cbs...)
}

func getValidateTimeEntryFn(c *api.Client) CallbackFn {
	if viper.GetBool(ALLOW_INCOMPLETE) {
		return nullCallback
	}

	return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
		return tei, validateTimeEntry(tei, c)
	}
}

func getPropsInteractiveFn(
	c *api.Client,
	dc suggestFn,
) CallbackFn {
	if viper.GetBool(INTERACTIVE) {
		return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			return askTimeEntryPropsInteractive(c, tei, dc)
		}
	}

	return nullCallback
}

func getDatesInteractiveFn() CallbackFn {
	if viper.GetBool(INTERACTIVE) {
		return func(tei dto.TimeEntryImpl) (dto.TimeEntryImpl, error) {
			return askTimeEntryDatesInteractive(tei)
		}
	}

	return nullCallback
}

func manageEntry(te dto.TimeEntryImpl, cbs ...CallbackFn) (
	dto.TimeEntryImpl, error) {
	return composeCallbacks(cbs...)(te)
}

func getErrorCode(err error) int {
	var e dto.Error
	if errors.As(err, &e) {
		return e.Code
	}

	return 0
}

func validateClosingTimeEntry(c *api.Client, workspace, userID string) error {
	te, err := c.GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
		Workspace: workspace,
		UserID:    userID,
	})

	if te == nil {
		return nil
	}

	if err != nil {
		return err
	}

	if err = validateTimeEntry(*te, c); err != nil {
		return fmt.Errorf("running time entry can't be ended: %w", err)
	}

	return nil
}

func out(te dto.TimeEntryImpl, c *api.Client) error {
	err := c.Out(api.OutParam{
		Workspace: te.WorkspaceID,
		UserID:    te.UserID,
		End:       te.TimeInterval.Start,
	})

	if getErrorCode(err) != 404 {
		return err
	}

	return nil
}

func createTimeEntry(te dto.TimeEntryImpl, c *api.Client) (
	dto.TimeEntryImpl, error) {
	return c.CreateTimeEntry(api.CreateTimeEntryParam{
		Workspace:   te.WorkspaceID,
		Billable:    te.Billable,
		Start:       te.TimeInterval.Start,
		End:         te.TimeInterval.End,
		ProjectID:   te.ProjectID,
		Description: te.Description,
		TagIDs:      te.TagIDs,
		TaskID:      te.TaskID,
	})

}

const noProject = "No Project"

func getProjectID(projectID string, w dto.Workspace, c *api.Client) (string, error) {
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

	projectID, err = ui.AskFromOptions("Choose your project:", projectsString, projectID)
	if err != nil || projectID == noProject || projectID == "" {
		return "", err
	}

	return strings.TrimSpace(projectID[0:strings.Index(projectID, " - ")]), nil
}

const noTask = "No Task"

func getTaskID(taskID, projectID string, w dto.Workspace, c *api.Client) (string, error) {
	tasks, err := c.GetTasks(api.GetTasksParam{
		Workspace:       w.ID,
		ProjectID:       projectID,
		PaginationParam: api.PaginationParam{AllPages: true},
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

func getDescription(description string, dc suggestFn) string {
	description, _ = ui.AskForText("Description:", description,
		ui.WithSuggestion(dc))
	return description
}

func getTagIDs(tagIDs []string, workspace string, c *api.Client) ([]string, error) {
	if len(tagIDs) > 0 && !viper.GetBool(INTERACTIVE) {
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

	var newTags []string
	if newTags, err = ui.AskManyFromOptions("Choose your tags:", tagsString, tagIDs); err != nil {
		return nil, nil
	}

	for i, t := range newTags {
		newTags[i] = strings.TrimSpace(t[0:strings.Index(t, " - ")])
	}

	return newTags, nil
}

func getUserId(c *api.Client) (string, error) {
	userId := viper.GetString(USER_ID)
	if len(userId) > 0 {
		return userId, nil
	}

	u, err := c.GetMe()
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

var noTimeEntryErr = errors.New("time entry was not found")

func getTimeEntry(
	id,
	workspace,
	userID string,
	c *api.Client,
) (dto.TimeEntryImpl, error) {
	id = strings.TrimSpace(strings.ToLower(id))

	mayNotFound := func(tei *dto.TimeEntryImpl, err error) (
		dto.TimeEntryImpl, error) {
		if err != nil {
			return dto.TimeEntryImpl{}, err
		}

		if tei == nil {
			return dto.TimeEntryImpl{}, noTimeEntryErr
		}

		return *tei, nil
	}

	switch id {
	case "^0", "current":
		id = "current"
	case "^1", "last":
		id = "last"
	}

	if id != "last" && id != "current" && !strings.HasPrefix(id, "^") {
		return mayNotFound(c.GetTimeEntry(api.GetTimeEntryParam{
			Workspace:   workspace,
			TimeEntryID: id,
		}))
	}

	if id == "current" {
		return mayNotFound(c.GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
			Workspace: workspace,
			UserID:    userID,
		}))
	}

	var b *bool
	if id == "latest" {
		b = nil
	}

	page := 1
	if strings.HasPrefix(id, "^") {
		var err error
		if page, err = strconv.Atoi(id[1:]); err != nil {
			return dto.TimeEntryImpl{}, fmt.Errorf(
				`n on "^n" must be a unsigned integer, you sent: %s`,
				id[1:],
			)
		}
	}

	list, err := c.GetUserTimeEntries(api.GetUserTimeEntriesParam{
		Workspace:      workspace,
		UserID:         userID,
		OnlyInProgress: b,
		PaginationParam: api.PaginationParam{
			PageSize: 1,
			Page:     page,
		},
	})

	if err != nil {
		return dto.TimeEntryImpl{}, err
	}

	if len(list) == 0 {
		return dto.TimeEntryImpl{}, noTimeEntryErr
	}

	return list[0], err
}

func addTimeEntryFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("not-billable", "n", false, "this time entry is not billable")
	cmd.Flags().String("task", "", "add a task to the entry")
	_ = completion.AddSuggestionsToFlag(cmd, "task", suggestWithClientAPI(suggestTasks))

	cmd.Flags().StringSliceP("tag", "T", []string{}, "add tags to the entry (can be used multiple times)")
	_ = completion.AddSuggestionsToFlag(cmd, "tag", suggestWithClientAPI(suggestTags))

	cmd.Flags().BoolP(ALLOW_INCOMPLETE, "A", false, "allow creation of incomplete time entries to be edited later (defaults to env $"+ENV_PREFIX+"_ALLOW_INCOMPLETE)")
	_ = viper.BindPFlag(ALLOW_INCOMPLETE, cmd.Flags().Lookup(ALLOW_INCOMPLETE))
	_ = viper.BindEnv(ALLOW_INCOMPLETE, ENV_PREFIX+"_ALLOW_INCOMPLETE")

	cmd.Flags().StringP("project", "p", "", "project to use for time entry")
	_ = completion.AddSuggestionsToFlag(cmd, "project", suggestWithClientAPI(suggestProjects))

	cmd.Flags().StringP("description", "d", "", "time entry description")
	_ = completion.AddSuggestionsToFlag(
		cmd,
		"description",
		suggestWithClientAPI(suggestDescription),
	)

	addPrintTimeEntriesFlags(cmd)

	// deprecations
	cmd.Flags().StringSlice("tags", []string{}, "add tags to the entry")
	_ = completion.AddSuggestionsToFlag(cmd, "tags", suggestWithClientAPI(suggestTags))
	_ = cmd.Flags().MarkDeprecated("tags", "use tag instead")
}

func addTimeEntryDateFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("when", "s", time.Now().Format(fullTimeFormat), "when the entry should be started, if not informed will use current time")
	cmd.Flags().StringP("when-to-close", "e", "", "when the entry should be closed, if not informed will let it open")
}

func fillTimeEntryWithFlags(tei dto.TimeEntryImpl, flags *pflag.FlagSet) (dto.TimeEntryImpl, error) {
	if flags.Changed("project") {
		tei.ProjectID, _ = flags.GetString("project")
	}

	if flags.Changed("description") {
		tei.Description, _ = flags.GetString("description")
	}

	if flags.Changed("task") {
		tei.TaskID, _ = flags.GetString("task")
	}

	if flags.Changed("tag") {
		tei.TagIDs, _ = flags.GetStringSlice("tag")
	}

	if flags.Changed("tags") {
		tei.TagIDs, _ = flags.GetStringSlice("tags")
	}

	if flags.Changed("not-billable") {
		b, _ := flags.GetBool("not-billable")
		tei.Billable = !b
	}

	var err error
	whenFlag := flags.Lookup("when")
	if whenFlag != nil && (whenFlag.Changed || whenFlag.DefValue != "") {
		whenString, _ := flags.GetString("when")
		var v time.Time
		if v, err = convertToTime(whenString); err != nil {
			return tei, err
		}
		tei.TimeInterval.Start = v
	}

	if flags.Changed("end-at") {
		whenString, _ := flags.GetString("end-at")
		var v time.Time
		if v, err = convertToTime(whenString); err != nil {
			return tei, err
		}
		tei.TimeInterval.End = &v
	}

	if flags.Changed("when-to-close") {
		whenString, _ := flags.GetString("when-to-close")
		var v time.Time
		if v, err = convertToTime(whenString); err != nil {
			return tei, err
		}
		tei.TimeInterval.End = &v
	}

	return tei, nil
}

func printMultipleTimeEntriesPreRun(cmd *cobra.Command, _ []string) error {
	viper.SetDefault(SHOW_TOTAL_DURATION, true)
	return viper.BindPFlag(SHOW_TOTAL_DURATION, cmd.Flags().Lookup("with-totals"))
}

func addPrintMultipleTimeEntriesFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("with-totals", "S", false, "add a totals line at the end")
}

func addPrintTimeEntriesFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "f", "",
		"golang text/template format to be applied on each time entry")
	cmd.Flags().BoolP("json", "j", false, "print as JSON")
	cmd.Flags().BoolP("csv", "v", false, "print as CSV")
	cmd.Flags().BoolP("quiet", "q", false, "print only ID")
	cmd.Flags().BoolP("md", "m", false, "print as Markdown")
	cmd.Flags().BoolP("duration-formatted", "D", false,
		"prints only the sum of duration formatted")
	cmd.Flags().BoolP("duration-float", "F", false,
		`prints only the sum of duration as a "float hour"`)
}

func addProjectFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("project", "p", "", "the name/id of the project to work on")
	cmd.MarkFlagRequired("project")

	_ = completion.AddSuggestionsToFlag(cmd, "project", suggestWithClientAPI(suggestProjects))
}

func getOpts(timeFormat string) []output.TimeEntryOutputOpt {
	opts := []output.TimeEntryOutputOpt{output.WithTimeFormat(timeFormat)}

	if viper.GetBool(SHOW_TASKS) {
		opts = append(opts, output.WithShowTasks())
	}

	if viper.GetBool(SHOW_TOTAL_DURATION) {
		opts = append(opts, output.WithTotalDuration())
	}

	return opts
}

func printTimeEntries(
	tes []dto.TimeEntry, cmd *cobra.Command, timeFormat string,
) error {
	var reportFn func(te []dto.TimeEntry, w io.Writer) error

	if b, _ := cmd.Flags().GetBool("md"); b {
		reportFn = output.TimeEntriesMarkdownPrint
	}

	if asJSON, _ := cmd.Flags().GetBool("json"); asJSON {
		reportFn = output.TimeEntriesJSONPrint
	}

	if asCSV, _ := cmd.Flags().GetBool("csv"); asCSV {
		reportFn = output.TimeEntriesCSVPrint
	}

	if format, _ := cmd.Flags().GetString("format"); format != "" {
		reportFn = output.TimeEntriesPrintWithTemplate(format)
	}

	if asQuiet, _ := cmd.Flags().GetBool("quiet"); asQuiet {
		reportFn = output.TimeEntriesPrintQuietly
	}

	if b, _ := cmd.Flags().GetBool("duration-float"); b {
		reportFn = output.TimeEntriesTotalDurationOnlyAsFloat
	}

	if b, _ := cmd.Flags().GetBool("duration-formatted"); b {
		reportFn = output.TimeEntriesTotalDurationOnlyFormatted
	}

	if reportFn == nil {
		reportFn = output.TimeEntriesPrint(getOpts(timeFormat)...)
	}

	return reportFn(tes, cmd.OutOrStdout())
}

func printTimeEntry(te *dto.TimeEntry, cmd *cobra.Command, timeFormat string) error {
	ts := make([]dto.TimeEntry, 0)
	if te != nil {
		ts = append(ts, *te)
	}

	return printTimeEntries(ts, cmd, timeFormat)
}
