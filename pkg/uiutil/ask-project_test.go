package uiutil_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/consoletest"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/lucassabreu/clockify-cli/pkg/uiutil"
	"github.com/stretchr/testify/assert"
)

func TestAskProjectShouldFail(t *testing.T) {
	tts := []struct {
		name  string
		param uiutil.AskProjectParam
		err   string
	}{
		{
			name:  "no ui",
			param: uiutil.AskProjectParam{},
			err:   "UI must be informed",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uiutil.AskProject(tt.param)
			if !assert.Error(t, err) {
				return
			}

			assert.Regexp(t, tt.err, err.Error())
		})
	}
}

var ps = []dto.Project{
	{ID: "p1", Name: "Project One"},
	{ID: "p2", Name: "Project Two", ClientID: "c1", ClientName: "Client One"},
	{ID: "p3", Name: "Project Tree"},
	{ID: "p4", Name: "Project Four"},
	{ID: "p5", Name: "Project Five"},
	{ID: "p6", Name: "Project Six"},
}

func TestAskProjectIsRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)
			ui.SetPageSize(10)

			p, err := uiutil.AskProject(uiutil.AskProjectParam{
				UI:        ui,
				ProjectID: "p2",
				Force:     true,
				Projects:  ps,
			})

			assert.Equal(t, &ps[3], p)
			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("project:")
			c.ExpectString("> p2")
			c.ExpectString("Client One")

			c.Send("four")
			c.ExpectString("four")
			c.ExpectString("> p4")
			c.SendLine()

			c.ExpectEOF()
		},
	)
}

func TestAskProjectIsntRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			p, err := uiutil.AskProject(uiutil.AskProjectParam{
				UI:        ui,
				Message:   "Which project?",
				ProjectID: "p2",
				Force:     false,
				Projects:  ps,
			})

			assert.Nil(t, p)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("project?")
			c.ExpectString("No Project")
			c.Send(string(terminal.KeyArrowUp))
			c.Send(string(terminal.KeyArrowUp))

			c.SendLine()

			c.ExpectEOF()
		},
	)
}

func TestAskProjectNoneSelected(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			p, err := uiutil.AskProject(uiutil.AskProjectParam{
				UI:        ui,
				Message:   "Which project?",
				ProjectID: "",
				Force:     false,
				Projects:  ps,
			})

			assert.Nil(t, p)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("project?")

			c.SendLine()

			c.ExpectEOF()
		},
	)
}
