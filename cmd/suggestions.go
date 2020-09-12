package cmd

import (
	"fmt"
	"strings"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/cmd/completion"
	"github.com/spf13/cobra"
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

func suggestTags(cmd *cobra.Command, args []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
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

func suggestProjects(cmd *cobra.Command, args []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
	projects, err := getProjects(c, "", false)
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

func suggestWorkspaces(cmd *cobra.Command, args []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
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

func suggestUsers(cmd *cobra.Command, args []string, toComplete string, c *api.Client) (completion.ValidArgs, error) {
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
