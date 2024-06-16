package uiutil

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/pkg/errors"
)

// AskProjectParam informs what options to display while asking for a project
type AskProjectParam struct {
	UI            ui.UI
	ProjectID     string
	Projects      []dto.Project
	ForceProjects bool
	Message       string
}

// NoProject is the text shown to not select a project
const NoProject = "No Project"

// AskProject asks the user for a project from options
func AskProject(p AskProjectParam) (*dto.Project, error) {
	if p.UI == nil {
		return nil, errors.New("UI must be informed")
	}

	if p.Projects == nil || len(p.Projects) == 0 {
		return nil, nil
	}
	if p.Message == "" {
		p.Message = "Choose your project:"
	}

	c, list := projectsToList(p.ProjectID, p.Projects)
	p.ProjectID = c

	if !p.ForceProjects {
		list = append([]string{NoProject}, list...)
	}

	id, err := p.UI.AskFromOptions(p.Message, list, p.ProjectID)
	if err != nil || id == NoProject || id == "" {
		return nil, err
	}

	id = strings.TrimSpace(id[0:strings.Index(id, " - ")])
	for i := range p.Projects {
		if p.Projects[i].ID == id {
			return &p.Projects[i], nil
		}
	}

	return nil, errors.New(`project with id "` + id + `" not found`)
}

func projectsToList(
	projectID string, projects []dto.Project) (string, []string) {
	list := make([]string, len(projects))
	found := -1
	nameSize := 0

	for i := range projects {
		list[i] = projects[i].ID + " - " + projects[i].Name
		if c := utf8.RuneCountInString(list[i]); nameSize < c {
			nameSize = c
		}

		if found == -1 && projects[i].ID == projectID {
			found = i
		}
	}

	format := fmt.Sprintf("%%-%ds| %%s", nameSize+1)
	for i := range projects {
		client := "Without Client"
		if projects[i].ClientID != "" {
			client = "Client: " + projects[i].ClientName +
				" (" + projects[i].ClientID + ")"
		}

		list[i] = fmt.Sprintf(format, list[i], client)
	}

	if found == -1 {
		return list[found], list
	}

	return "", list
}
