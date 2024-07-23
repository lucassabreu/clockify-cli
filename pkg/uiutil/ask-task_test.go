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

func TestAskTaskShouldFail(t *testing.T) {
	tts := []struct {
		name  string
		param uiutil.AskTaskParam
		err   string
	}{
		{
			name:  "no ui",
			param: uiutil.AskTaskParam{},
			err:   "UI must be informed",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uiutil.AskTask(tt.param)
			if !assert.Error(t, err) {
				return
			}

			assert.Regexp(t, tt.err, err.Error())
		})
	}
}

var tks = []dto.Task{
	{ID: "t1", Name: "Task One"},
	{ID: "t2", Name: "Task Two"},
	{ID: "t3", Name: "Task Tree"},
	{ID: "t4", Name: "Task Four"},
	{ID: "t5", Name: "Task Five"},
	{ID: "t6", Name: "Task Six"},
}

func TestAskTaskIsRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)
			ui.SetPageSize(10)

			p, err := uiutil.AskTask(uiutil.AskTaskParam{
				UI:     ui,
				TaskID: "t2",
				Force:  true,
				Tasks:  tks,
			})

			assert.Equal(t, &tks[3], p)
			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("task:")
			c.ExpectString("> t2")

			c.Send("four")
			c.ExpectString("four")
			c.ExpectString("> t4")
			c.SendLine()

			c.ExpectEOF()
		},
	)
}

func TestAskTaskIsntRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			p, err := uiutil.AskTask(uiutil.AskTaskParam{
				UI:      ui,
				Message: "Which task?",
				TaskID:  "t2",
				Force:   false,
				Tasks:   tks,
			})

			assert.Nil(t, p)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("task?")
			c.ExpectString("No Task")
			c.Send(string(terminal.KeyArrowUp))
			c.Send(string(terminal.KeyArrowUp))

			c.SendLine()

			c.ExpectEOF()
		},
	)
}

func TestAskTaskNoneSelected(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			p, err := uiutil.AskTask(uiutil.AskTaskParam{
				UI:      ui,
				Message: "Which task?",
				TaskID:  "",
				Force:   false,
				Tasks:   tks,
			})

			assert.Nil(t, p)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("task?")
			c.ExpectString("No Task")

			c.SendLine()

			c.ExpectEOF()
		},
	)
}
