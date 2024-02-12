package uiutil

import (
	"strings"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/pkg/errors"
)

// AskTaskParam informs what options to display while asking for a task
type AskTaskParam struct {
	UI      ui.UI
	TaskID  string
	Tasks   []dto.Task
	Force   bool
	Message string
}

const noTask = "No Task"

// AskTask asks the user for a task from options
func AskTask(p AskTaskParam) (*dto.Task, error) {
	if p.UI == nil {
		return nil, errors.New("UI must be informed")
	}

	if p.Tasks == nil || len(p.Tasks) == 0 {
		return nil, nil
	}
	if p.Message == "" {
		p.Message = "Choose your task:"
	}

	list := make([]string, len(p.Tasks))
	found := -1

	for i := range p.Tasks {
		list[i] = p.Tasks[i].ID + " - " + p.Tasks[i].Name
		if found == -1 && p.Tasks[i].ID == p.TaskID {
			found = i
		}
	}

	if found == -1 {
		p.TaskID = ""
	} else {
		p.TaskID = list[found]
	}

	if !p.Force {
		list = append([]string{noTask}, list...)
	}

	id, err := p.UI.AskFromOptions(p.Message, list, p.TaskID)
	if err != nil || id == noTask || id == "" {
		return nil, err
	}

	id = strings.TrimSpace(id[0:strings.Index(id, " - ")])
	for i := range p.Tasks {
		if p.Tasks[i].ID == id {
			return &p.Tasks[i], nil
		}
	}

	return nil, errors.New(`task with id "` + id + `" not found`)
}
