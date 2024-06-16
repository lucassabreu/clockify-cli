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

func TestAskTagsShouldFail(t *testing.T) {
	tts := []struct {
		name  string
		param uiutil.AskTagsParam
		err   string
	}{
		{
			name:  "no ui",
			param: uiutil.AskTagsParam{},
			err:   "UI must be informed",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uiutil.AskTags(tt.param)
			if !assert.Error(t, err) {
				return
			}

			assert.Regexp(t, tt.err, err.Error())
		})
	}
}

var tags = []dto.Tag{
	{ID: "t1", Name: "Tag One"},
	{ID: "t2", Name: "Tag Two"},
	{ID: "t3", Name: "Tag Tree"},
	{ID: "t4", Name: "Tag Four"},
	{ID: "t5", Name: "Tag Five"},
	{ID: "t6", Name: "Tag Six"},
}

func TestAskTags(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)
			ui.SetPageSize(10)

			ts, err := uiutil.AskTags(uiutil.AskTagsParam{
				UI:     ui,
				TagIDs: []string{"t2", "t4"},
				Force:  true,
				Tags:   tags,
			})

			if !assert.Equal(
				t,
				[]dto.Tag{
					{ID: "t1", Name: "Tag One"},
					{ID: "t2", Name: "Tag Two"},
					{ID: "t5", Name: "Tag Five"},
					{ID: "t6", Name: "Tag Six"},
				},
				ts,
			) {
				return nil
			}

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("tags:")
			c.ExpectString("[x]")
			c.ExpectString("[x]")

			c.Send("one ")
			c.Send("four ")

			c.Send("f")
			c.Send(string(terminal.KeyArrowDown))
			c.Send(" ")

			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))

			c.SendLine(" ")

			c.ExpectEOF()
		},
	)
}

func TestAskTagsIsRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			ts, err := uiutil.AskTags(uiutil.AskTagsParam{
				UI:      ui,
				Message: "Which tags?",
				TagIDs:  []string{"t2"},
				Force:   true,
				Tags:    tags,
			})

			assert.Equal(
				t,
				[]dto.Tag{
					{ID: "t1", Name: "Tag One"},
				},
				ts,
			)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("tags?")
			c.ExpectString("[x]")

			c.SendLine(string(terminal.KeyArrowLeft))

			c.ExpectString("at least one")
			c.Send(string(terminal.KeyArrowLeft))

			c.SendLine(" ")

			c.ExpectEOF()
		},
	)
}

func TestAskTagsIsntRequired(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			ts, err := uiutil.AskTags(uiutil.AskTagsParam{
				UI:      ui,
				Message: "Which tags?",
				TagIDs:  []string{"t2"},
				Force:   false,
				Tags:    tags,
			})

			assert.Equal(
				t,
				[]dto.Tag{},
				ts,
			)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("tags?")
			c.ExpectString("[x]")

			c.SendLine(string(terminal.KeyArrowLeft))

			c.ExpectEOF()
		},
	)
}

func TestAskTagsNoneSelected(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui := ui.NewUI(in, out, out)

			ts, err := uiutil.AskTags(uiutil.AskTagsParam{
				UI:      ui,
				Message: "Which tags?",
				TagIDs:  []string{},
				Force:   false,
				Tags:    tags,
			})

			assert.Equal(
				t,
				[]dto.Tag{},
				ts,
			)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("tags?")
			c.SendLine()
			c.ExpectEOF()
		},
	)
}
