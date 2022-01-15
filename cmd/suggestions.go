package cmd

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func suggestWithClientAPI(
	fn func(cmd *cobra.Command, args []string, toComplete string, c *api.Client) (completion.ValidArgs, error),
) func(cmd *cobra.Command, args []string, toComplete string) (completion.ValidArgs, error) {
	return func(cmd *cobra.Command, args []string, toComplete string) (completion.ValidArgs, error) {
		c, err := getAPIClient()
		if err != nil {
			return completion.EmptyValidArgs(), err
		}

		return fn(cmd, args, toComplete, c)
	}
}

func suggestTags(_ *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	tags, err := getTags(c, "", false)
	if err != nil {
		return completion.EmptyValidArgs(), err
	}

	va := make(completion.ValigsArgsMap)
	toComplete = strings.ToLower(toComplete)
	for _, tag := range tags {
		if toComplete != "" && !strings.Contains(tag.ID, toComplete) {
			continue
		}
		va.Set(tag.ID, tag.Name)
	}

	return va, nil
}

func suggestTasks(cmd *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return completion.EmptyValidArgs(), err
	}

	w := viper.GetString(WORKSPACE)
	if viper.GetBool(ALLOW_NAME_FOR_ID) {
		project, err = getProjectByNameOrId(c, w, project)
		if err != nil {
			return completion.EmptyValidArgs(), err
		}
	}

	tasks, err := c.GetTasks(api.GetTasksParam{
		Workspace: w,
		ProjectID: project,
	})

	if err != nil || len(tasks) == 0 {
		return completion.EmptyValidArgs(), err
	}

	va := make(completion.ValigsArgsMap)
	toComplete = strings.ToLower(toComplete)
	for _, task := range tasks {
		if toComplete != "" && !strings.Contains(task.ID, toComplete) {
			continue
		}
		va.Set(task.ID, task.Name)
	}

	return va, nil
}

func suggestProjects(_ *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	b := false
	projects, err := c.GetProjects(api.GetProjectsParam{
		Workspace:       viper.GetString(WORKSPACE),
		Archived:        &b,
		PaginationParam: api.AllPages(),
	})
	if err != nil {
		return completion.EmptyValidArgs(), err
	}

	va := make(completion.ValigsArgsMap)
	toComplete = strings.ToLower(toComplete)
	for _, project := range projects {
		if toComplete != "" && !strings.Contains(project.ID, toComplete) {
			continue
		}
		va.Set(project.ID, project.Name)
	}

	return va, nil
}

func suggestWorkspaces(_ *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	workspaces, err := getWorkspaces(c, "")

	if err != nil {
		return completion.EmptyValidArgs(), err
	}

	va := make(completion.ValigsArgsMap)
	toComplete = strings.ToLower(toComplete)
	for _, workspace := range workspaces {
		if toComplete != "" && !strings.Contains(workspace.ID, toComplete) {
			continue
		}
		va.Set(workspace.ID, workspace.Name)
	}

	return va, nil
}

func suggestUsers(_ *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	users, err := getUsers(c, "")

	if err != nil {
		return completion.EmptyValidArgs(), err
	}

	va := make(completion.ValigsArgsMap)
	toComplete = strings.ToLower(toComplete)
	for _, user := range users {
		if toComplete != "" && !strings.Contains(user.ID, toComplete) {
			continue
		}
		va.Set(user.ID, fmt.Sprintf("%s (%s)", user.Name, user.Email))
	}

	return va, nil
}

func suggestDescription(_ *cobra.Command, _ []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	if viper.GetBool(DESCR_AUTOCOMP) {
		return completion.EmptyValidArgs(), nil
	}

	dc := newDescriptionCompleter(
		c,
		viper.GetString(WORKSPACE),
		viper.GetString(USER_ID),
		viper.GetInt(DESCR_AUTOCOMP_DAYS),
	)

	return completion.ValigsArgsSlide(dc.suggestFn(toComplete)), nil
}
