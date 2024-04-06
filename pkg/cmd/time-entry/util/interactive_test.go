package util

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/consoletest"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/lucassabreu/clockify-cli/pkg/uiutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPropsInteractive_ShouldSkip_WhenDisabled(t *testing.T) {
	f := mocks.NewMockFactory(t)
	f.EXPECT().Config().Return(&mocks.SimpleConfig{
		Interactive: false,
	})
	s := GetPropsInteractiveFn(nil, f)

	te := TimeEntryDTO{}
	te2, err := s(te)

	assert.NoError(t, err)
	assert.Equal(t, te, te2)
}

func TestGetPropsInteractive_ShouldAskValues(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			c := mocks.NewMockClient(t)

			c.EXPECT().GetWorkspace(mock.Anything).
				Return(dto.Workspace{ID: "w"}, nil)

			c.EXPECT().GetProjects(mock.Anything).
				Return(
					[]dto.Project{
						{ID: "1", Name: "First"},
						{ID: "2", Name: "Second",
							ClientID: "1", ClientName: "Client One"},
						{ID: "3", Name: "Third",
							ClientID: "2", ClientName: "Client Two"},
						{ID: "4", Name: "Fourth"},
					},
					nil,
				)

			c.EXPECT().GetTasks(mock.Anything).
				Return(
					[]dto.Task{
						{ID: "t1", Name: "First"},
						{ID: "t2", Name: "Second"},
						{ID: "t3", Name: "Third"},
					},
					nil,
				)

			c.EXPECT().GetTags(mock.Anything).
				Return(
					[]dto.Tag{
						{ID: "tag1", Name: "meeting"},
						{ID: "tag2", Name: "backend"},
						{ID: "tag3", Name: "frontend"},
					},
					nil,
				)

			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			te, err := GetPropsInteractiveFn(
				func(string) []string { return []string{} },
				f,
			)(TimeEntryDTO{
				Workspace: "w",
			})

			assert.NoError(t, err)
			assert.Equal(t,
				TimeEntryDTO{
					Workspace:   "w",
					ProjectID:   "3",
					TaskID:      "t2",
					Description: "a unique description",
					TagIDs:      []string{"tag2", "tag3"},
				},
				te,
			)

			return err
		}, func(c consoletest.ExpectConsole) {
			c.ExpectString("Choose your project:")
			c.ExpectString(uiutil.NoProject)
			c.ExpectString("1 - First  | Without Client")
			c.ExpectString("2 - Second | Client: Client One (1)")
			c.ExpectString("3 - Third  | Client: Client Two (2)")
			c.ExpectString("4 - Fourth | Without Client")

			c.Send("ir")
			c.ExpectString("1 - First  | Without Client")
			c.ExpectString("3 - Third  | Client: Client Two (2)")

			c.Send(string(terminal.KeyArrowDown))
			c.SendLine()

			c.ExpectString("Choose your task:")
			c.ExpectString(uiutil.NoTask)
			c.ExpectString("t1 - First")
			c.ExpectString("t2 - Second")
			c.ExpectString("t3 - Third")

			c.SendLine("2")

			c.ExpectString("Description:")
			c.SendLine("a unique description")

			c.ExpectString("Choose your tags:")

			c.ExpectString("tag1 - meeting")
			c.ExpectString("tag2 - backend")
			c.ExpectString("tag3 - frontend")

			c.Send("end")
			c.Send(string(terminal.KeyArrowRight))
			c.SendLine()

			c.ExpectEOF()
		})
}

func TestGetPropsInteractive_ShouldAllowEmptyValues(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			c := mocks.NewMockClient(t)

			c.EXPECT().GetWorkspace(mock.Anything).
				Return(dto.Workspace{ID: "w"}, nil)

			c.EXPECT().GetProjects(mock.Anything).
				Return(
					[]dto.Project{
						{ID: "1", Name: "First"},
					},
					nil,
				)

			c.EXPECT().GetTags(mock.Anything).
				Return(
					[]dto.Tag{
						{ID: "tag1", Name: "meeting"},
						{ID: "tag2", Name: "backend"},
						{ID: "tag3", Name: "frontend"},
					},
					nil,
				)

			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			te, err := GetPropsInteractiveFn(
				func(string) []string { return []string{} },
				f,
			)(TimeEntryDTO{
				Workspace: "w",
			})

			assert.NoError(t, err)
			assert.Equal(t,
				TimeEntryDTO{
					Workspace:   "w",
					ProjectID:   "",
					TaskID:      "",
					Description: "",
					TagIDs:      nil,
				},
				te,
			)

			return err
		}, func(c consoletest.ExpectConsole) {
			c.ExpectString("Choose your project:")
			c.ExpectString(uiutil.NoProject)
			c.ExpectString("1 - First | Without Client")

			c.SendLine("")

			c.ExpectString("Description:")
			c.SendLine("")

			c.ExpectString("Choose your tags:")

			c.ExpectString("tag1 - meeting")
			c.ExpectString("tag2 - backend")
			c.ExpectString("tag3 - frontend")

			c.Send("end")
			c.SendLine("")

			c.ExpectEOF()
		})
}

func TestGetPropsInteractive_ShouldUseInputAsSelected(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			c := mocks.NewMockClient(t)

			c.EXPECT().GetWorkspace(mock.Anything).
				Return(dto.Workspace{ID: "w"}, nil)

			c.EXPECT().GetProjects(mock.Anything).
				Return(
					[]dto.Project{
						{ID: "1", Name: "First"},
						{ID: "2", Name: "Second",
							ClientID: "1", ClientName: "Client One"},
						{ID: "3", Name: "Third",
							ClientID: "2", ClientName: "Client Two"},
						{ID: "4", Name: "Fourth"},
					},
					nil,
				)

			c.EXPECT().GetTasks(mock.Anything).
				Return(
					[]dto.Task{
						{ID: "t1", Name: "First"},
						{ID: "t2", Name: "Second"},
						{ID: "t3", Name: "Third"},
					},
					nil,
				)

			c.EXPECT().GetTags(mock.Anything).
				Return(
					[]dto.Tag{
						{ID: "tag1", Name: "meeting"},
						{ID: "tag2", Name: "backend"},
						{ID: "tag3", Name: "frontend"},
					},
					nil,
				)

			input := TimeEntryDTO{
				Workspace:   "w",
				ProjectID:   "3",
				TaskID:      "t2",
				Description: "a unique description",
				TagIDs:      []string{"tag2", "tag3"},
			}

			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			output, err := GetPropsInteractiveFn(
				func(string) []string { return []string{} },
				f,
			)(input)

			assert.NoError(t, err)
			assert.Equal(t, output, input)

			return err
		}, func(c consoletest.ExpectConsole) {
			c.ExpectString("Choose your project:")
			c.ExpectString("Third")
			c.SendLine("")

			c.ExpectString("Choose your task:")
			c.SendLine("")

			c.ExpectString("Description:")
			c.ExpectString("a unique description")
			c.SendLine("")

			c.ExpectString("Choose your tags:")
			c.SendLine("")

			c.ExpectEOF()
		})
}

func TestGetPropsInteractive_ShouldForceAnswer_WhenWorkspaceForces(
	t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			c := mocks.NewMockClient(t)

			c.EXPECT().GetWorkspace(mock.Anything).
				Return(
					dto.Workspace{
						ID: "w",
						Settings: dto.WorkspaceSettings{
							ForceProjects:    true,
							ForceTasks:       true,
							ForceDescription: true,
							ForceTags:        true,
						},
					},
					nil,
				)

			c.EXPECT().GetProjects(mock.Anything).
				Return(
					[]dto.Project{
						{ID: "1", Name: "First"},
						{ID: "2", Name: "Second",
							ClientID: "1", ClientName: "Client One"},
					},
					nil,
				)

			c.EXPECT().GetTasks(mock.Anything).
				Return(
					[]dto.Task{
						{ID: "t1", Name: "First"},
						{ID: "t2", Name: "Second"},
					},
					nil,
				)

			c.EXPECT().GetTags(mock.Anything).
				Return(
					[]dto.Tag{
						{ID: "tag1", Name: "meeting"},
						{ID: "tag2", Name: "backend"},
						{ID: "tag3", Name: "frontend"},
					},
					nil,
				)

			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			output, err := GetPropsInteractiveFn(
				func(string) []string { return []string{} },
				f,
			)(TimeEntryDTO{Workspace: "w"})

			assert.NoError(t, err)
			assert.Equal(t,
				TimeEntryDTO{
					Workspace:   "w",
					ProjectID:   "1",
					TaskID:      "t1",
					Description: "something",
					TagIDs:      []string{"tag1"},
				},
				output,
			)

			return err
		}, func(c consoletest.ExpectConsole) {
			c.ExpectString("Choose your project:")
			c.SendLine("")

			c.ExpectString("Choose your task:")
			c.SendLine("")

			c.ExpectString("Description:")
			c.SendLine("")
			c.ExpectString("description should be informed")
			c.SendLine("something")

			c.ExpectString("Choose your tags:")
			c.SendLine("")
			c.ExpectString("at least one tag should be selected")
			c.SendLine(" ")

			c.ExpectEOF()
		})
}

func TestGetPropsInteractive_ShouldNotAsk_WhenThereAreNoOptions(
	t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			c := mocks.NewMockClient(t)

			c.EXPECT().GetWorkspace(mock.Anything).
				Return(dto.Workspace{ID: "w"}, nil)

			c.EXPECT().GetProjects(mock.Anything).
				Return([]dto.Project{}, nil)

			c.EXPECT().GetTags(mock.Anything).
				Return([]dto.Tag{}, nil)

			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Client().Return(c, nil)
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			output, err := GetPropsInteractiveFn(
				func(string) []string { return []string{} },
				f,
			)(TimeEntryDTO{Workspace: "w"})

			assert.NoError(t, err)
			assert.Equal(t,
				TimeEntryDTO{
					Workspace:   "w",
					ProjectID:   "",
					TaskID:      "",
					Description: "something",
					TagIDs:      nil,
				},
				output,
			)

			return err
		}, func(c consoletest.ExpectConsole) {
			c.ExpectString("Description:")
			c.SendLine("something")

			c.ExpectEOF()
		})
}

func TestGetDatesInteractive_ShouldSkip_WhenDisabled(t *testing.T) {
	f := mocks.NewMockFactory(t)
	f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: false})

	s := GetDatesInteractiveFn(f)

	te := TimeEntryDTO{}
	te2, err := s(te)

	assert.NoError(t, err)
	assert.Equal(t, te, te2)
}

func TestGetDatesInteractive_ShouldValidateString_WhenWrongFormat(
	t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			f := mocks.NewMockFactory(t)
			f.EXPECT().UI().Return(ui.NewUI(in, out, out))
			f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

			s := GetDatesInteractiveFn(f)

			te, err := s(TimeEntryDTO{
				Start: timehlp.Now(),
			})

			assert.NoError(t, err)
			start, _ := timehlp.ConvertToTime(timehlp.FullTimeFormat)
			assert.Equal(t,
				TimeEntryDTO{
					Start: start,
					End:   nil,
				},
				te)

			return nil
		},
		func(c consoletest.ExpectConsole) {
			wrongInputs := func() {
				for _, v := range []string{
					"wrong",
					"99:99",
					"99:99:99",
				} {
					c.SendLine(v)
					c.ExpectString("Sorry, your reply was invalid")
				}
			}

			c.ExpectString("Start:")
			wrongInputs()
			c.SendLine(timehlp.FullTimeFormat)

			c.ExpectString("End")
			wrongInputs()
			c.SendLine("")

			c.ExpectEOF()
		},
	)
}

func TestGetDatesInteractive_ShouldAccept_ValideTimeFormats(t *testing.T) {
	// toTimeRef := func(t time.Time) *time.Time { return &t }
	fromTimeString := func(s string) TimeEntryDTO {
		t, _ := timehlp.ConvertToTime(s)
		return TimeEntryDTO{
			Start: t,
			End:   &t,
		}
	}

	tts := []struct {
		timeString string
		input      TimeEntryDTO
		output     TimeEntryDTO
	}{
		{
			timeString: timehlp.FullTimeFormat,
			output:     fromTimeString(timehlp.FullTimeFormat),
		},
		{
			timeString: timehlp.SimplerTimeFormat,
			output:     fromTimeString(timehlp.SimplerTimeFormat),
		},
		{
			timeString: timehlp.OnlyTimeFormat,
			output:     fromTimeString(timehlp.OnlyTimeFormat),
		},
		{
			timeString: timehlp.SimplerOnlyTimeFormat,
			output:     fromTimeString(timehlp.SimplerOnlyTimeFormat),
		},
		{
			timeString: timehlp.NowTimeFormat,
			output:     fromTimeString(timehlp.NowTimeFormat),
		},
		{
			timeString: "",
			input:      fromTimeString(timehlp.FullTimeFormat),
			output:     fromTimeString(timehlp.FullTimeFormat),
		},
	}

	for i := range tts {
		tt := &tts[i]

		t.Run(tt.timeString, func(t *testing.T) {
			consoletest.RunTestConsole(t,
				func(out consoletest.FileWriter, in consoletest.FileReader) error {
					f := mocks.NewMockFactory(t)
					f.EXPECT().UI().Return(ui.NewUI(in, out, out))
					f.EXPECT().Config().Return(&mocks.SimpleConfig{Interactive: true})

					s := GetDatesInteractiveFn(f)

					te, err := s(tt.input)

					assert.NoError(t, err)
					assert.Equal(t, tt.output, te)
					return nil
				},
				func(c consoletest.ExpectConsole) {
					c.ExpectString("Start:")
					c.SendLine(tt.timeString)

					c.ExpectString("End")
					c.SendLine(tt.timeString)

					c.ExpectEOF()
				},
			)
		})
	}
}
