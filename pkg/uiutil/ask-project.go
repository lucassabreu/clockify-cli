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

const noProject = "No Project"

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

	list := make([]string, len(p.Projects))
	found := -1
	nameSize := 0

	for i := range p.Projects {
		list[i] = p.Projects[i].ID + " - " + p.Projects[i].Name
		if c := utf8.RuneCountInString(list[i]); nameSize < c {
			nameSize = c
		}

		if found == -1 && p.Projects[i].ID == p.ProjectID {
			found = i
		}
	}

	format := fmt.Sprintf("%%-%ds| %%s", nameSize+1)
	for i := range p.Projects {
		client := "Without Client"
		if p.Projects[i].ClientID != "" {
			client = "Client: " + p.Projects[i].ClientName +
				" (" + p.Projects[i].ClientName + ")"
		}

		list[i] = fmt.Sprintf(format, list[i], client)
	}

	if found == -1 {
		p.ProjectID = ""
	} else {
		p.ProjectID = list[found]
	}

	if !p.ForceProjects {
		list = append([]string{noProject}, list...)
	}

	id, err := p.UI.AskFromOptions(p.Message, list, p.ProjectID)
	if err != nil || id == noProject || id == "" {
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
