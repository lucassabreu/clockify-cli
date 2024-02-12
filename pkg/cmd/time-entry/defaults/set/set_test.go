package set_test

import (
	"errors"
	"io"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/consoletest"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/set"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	. "github.com/lucassabreu/clockify-cli/pkg/output/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var bTrue = true
var bFalse = false

func TestNewCmdSet_ShouldAskInfo_WhenInteractive(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			f := mocks.NewMockFactory(t)

			f.EXPECT().UI().Return(ui.NewUI(in, out, out))

			ted := mocks.NewMockTimeEntryDefaults(t)
			ted.EXPECT().Read().Return(defaults.DefaultTimeEntry{}, nil)
			ted.On("Write", mock.Anything).Return(nil)
			f.EXPECT().TimeEntryDefaults().Return(ted)

			f.EXPECT().GetWorkspaceID().Return("w", nil)

			c := mocks.NewMockClient(t)

			c.EXPECT().GetProjects(api.GetProjectsParam{
				Workspace:       "w",
				Archived:        &bFalse,
				PaginationParam: api.AllPages(),
			}).
				Return([]dto.Project{
					{ID: "p1", Name: "first"},
					{ID: "p2", Name: "second",
						ClientID: "c", ClientName: "Myself"},
					{ID: "p3", Name: "third"},
				}, nil)

			c.EXPECT().GetTasks(api.GetTasksParam{
				Workspace:       "w",
				ProjectID:       "p3",
				Active:          true,
				PaginationParam: api.AllPages(),
			}).
				Return([]dto.Task{
					{ID: "t", Name: "first"},
					{ID: "t2", Name: "second"},
					{ID: "t3", Name: "third"},
				}, nil)

			c.EXPECT().GetTags(api.GetTagsParam{
				Workspace:       "w",
				Archived:        &bFalse,
				PaginationParam: api.AllPages(),
			}).
				Return([]dto.Tag{
					{ID: "tg1", Name: "first"},
					{ID: "tg2", Name: "second"},
					{ID: "tg3", Name: "third"},
				}, nil)

			f.EXPECT().Client().Return(c, nil)

			f.EXPECT().Config().Return(&mocks.SimpleConfig{
				AllowNameForID:            true,
				Interactive:               true,
				InteractivePageSizeNumber: 7,
			})

			var d defaults.DefaultTimeEntry
			cmd := set.NewCmdSet(f, func(_ OutputFlags, _ io.Writer,
				dte defaults.DefaultTimeEntry) error {
				d = dte
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			cmd.SetArgs([]string{
				"-p=not found",
				"--task=nf",
				"-T=nf",
				"--not-billable",
			})
			_, err := cmd.ExecuteC()

			if !assert.NoError(t, err) {
				return err
			}

			assert.Equal(t, "w", d.Workspace)
			assert.Equal(t, "p3", d.ProjectID)
			assert.Equal(t, "t2", d.TaskID)
			assert.Equal(t, []string{"tg1", "tg2", "tg3"}, d.TagIDs)
			assert.Equal(t, &bTrue, d.Billable)

			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("Choose your project:")
			c.ExpectString("> No Project")
			c.ExpectString("first  | Without Client")
			c.ExpectString("second | Client: Myself")
			c.ExpectString("third  | Without Client")

			c.Send("without")
			c.SendLine(string(terminal.KeyArrowDown))

			c.ExpectString("Choose your task:")
			c.ExpectString("> No Task")
			c.ExpectString("first")
			c.ExpectString("second")
			c.ExpectString("third")
			c.SendLine("sec")

			c.ExpectString("Choose your tags:")
			c.ExpectString("first")
			c.ExpectString("second")
			c.ExpectString("third")
			c.SendLine(string(terminal.KeyArrowRight))

			c.ExpectString("Should be billable?")
			c.SendLine("y")

			c.ExpectEOF()
		},
	)
}

func runCmd(f cmdutil.Factory, args []string) (
	d defaults.DefaultTimeEntry, reported bool, err error) {

	cmd := set.NewCmdSet(f, func(_ OutputFlags, _ io.Writer,
		dte defaults.DefaultTimeEntry) error {
		reported = true
		d = dte
		return nil
	})

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.SetArgs(args)
	_, err = cmd.ExecuteC()

	return d, reported, err
}

func TestNewCmdSet_ShouldFail_WhenInvalidArgs(t *testing.T) {
	tts := []struct {
		name    string
		args    []string
		err     string
		factory func(t *testing.T) cmdutil.Factory
	}{
		{
			name: "can't be not billable and billable",
			args: []string{"--billable", "--not-billable"},
			err:  ".*flags can't be used together.*",
			factory: func(*testing.T) cmdutil.Factory {
				return mocks.NewMockFactory(t)
			},
		},
		{
			name: "can't read file",
			err:  "failed",
			factory: func(*testing.T) cmdutil.Factory {
				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					errors.New("failed"),
				)

				f := mocks.NewMockFactory(t)
				f.EXPECT().TimeEntryDefaults().Return(ted)

				return f
			},
		},
		{
			name: "failed to get client",
			args: []string{"--project", "p1"},
			err:  "failed",
			factory: func(*testing.T) cmdutil.Factory {
				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f := mocks.NewMockFactory(t)
				f.EXPECT().TimeEntryDefaults().Return(ted)
				f.EXPECT().GetWorkspaceID().Return("w", nil)
				f.EXPECT().Client().Return(
					mocks.NewMockClient(t),
					errors.New("failed"),
				)

				return f
			},
		},
		{
			name: "can't get workspace",
			err:  "failed",
			factory: func(*testing.T) cmdutil.Factory {
				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f := mocks.NewMockFactory(t)
				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("", errors.New("failed"))

				return f
			},
		},
		{
			name: "can't get project",
			err:  "failed",
			args: []string{"--project", "p"},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)
				f.EXPECT().GetWorkspaceID().Return("w", nil)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p",
					Hydrate:   false,
				}).Return(nil, errors.New("failed"))

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't find task",
			err:  `can't find task with ID "tk" on project "p"`,
			args: []string{
				"--project", "p",
				"--task=tk",
			},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p",
					Hydrate:   true,
				}).Return(&dto.Project{ID: "p", Name: "project"}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't find tag",
			err:  "can't find tag with ID \"tg\"",
			args: []string{
				"--project", "p",
				"-T", "tg",
			},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProject(api.GetProjectParam{
					Workspace: "w",
					ProjectID: "p",
					Hydrate:   false,
				}).Return(&dto.Project{ID: "p", Name: "project"}, nil)

				cl.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return([]dto.Tag{{ID: "not that"}}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't look for tag",
			err:  "failed",
			args: []string{
				"-T", "tg",
			},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
				})

				cl := mocks.NewMockClient(t)

				cl.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).Return(nil, errors.New("failed"))

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't find project by name",
			err:  "No project with id or name containing 'p' was found",
			args: []string{"--project", "p"},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProjects(mock.Anything).
					Return([]dto.Project{}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't find task by name",
			err:  "No task with id or name containing 'task' was found",
			args: []string{"--project", "project", "--task=task"},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProjects(mock.Anything).
					Return([]dto.Project{{ID: "p", Name: "project"}}, nil)

				cl.EXPECT().GetTasks(api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p",
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Task{{ID: "tk", Name: "other"}}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't find tag by name",
			err:  "No tag with id or name containing 'tag' was found",
			args: []string{
				"--project", "project",
				"--task=task",
				"-T=tag",
			},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

				cl := mocks.NewMockClient(t)
				cl.EXPECT().GetProjects(mock.Anything).
					Return([]dto.Project{{ID: "p", Name: "project"}}, nil)

				cl.EXPECT().GetTasks(api.GetTasksParam{
					Workspace:       "w",
					ProjectID:       "p",
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Task{{ID: "tk", Name: "task"}}, nil)

				cl.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Tag{{ID: "tg", Name: "other"}}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
		{
			name: "can't set task without project",
			err:  "can't set task without project",
			args: []string{"--task=task"},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				return f
			},
		},
		{
			name: "can't find tag by name (no project)",
			err:  "No tag with id or name containing 'tag2' was found",
			args: []string{"-T=tag", "-T=tag2"},
			factory: func(*testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)

				ted := mocks.NewMockTimeEntryDefaults(t)
				ted.EXPECT().Read().Return(
					defaults.DefaultTimeEntry{},
					defaults.DefaultsFileNotFoundErr,
				)

				f.EXPECT().TimeEntryDefaults().Return(ted)

				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

				cl := mocks.NewMockClient(t)

				cl.EXPECT().GetTags(api.GetTagsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Tag{{ID: "tag", Name: "other"}}, nil)

				f.EXPECT().Client().Return(cl, nil)

				return f
			},
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			_, called, err := runCmd(tt.factory(t), tt.args)
			if !assert.Error(t, err, "should have failed") {
				return
			}
			assert.False(t, called)
			assert.Regexp(t, tt.err, err)
		})
	}
}

func TestNewCmdSet_ShouldUpdateDefaultsFile_OnlyByFlags(t *testing.T) {
	tts := []struct {
		name     string
		args     []string
		current  defaults.DefaultTimeEntry
		expected defaults.DefaultTimeEntry
	}{
		{
			name: "no arguments, no changes",
			args: []string{},
			current: defaults.DefaultTimeEntry{
				Workspace: "w1", ProjectID: "p1"},
			expected: defaults.DefaultTimeEntry{
				Workspace: "w1", ProjectID: "p1"},
		},
		{
			name: "all arguments",
			args: []string{
				"-p=p2",
				"--task=t2",
				"-T=tg1", "-T=tg2",
				"--billable",
			},
			expected: defaults.DefaultTimeEntry{
				Workspace: "w2",
				ProjectID: "p2",
				TaskID:    "t2",
				Billable:  &bTrue,
				TagIDs:    []string{"tg1", "tg2"},
			},
		},
		{
			name: "not billable",
			args: []string{"--not-billable"},
			current: defaults.DefaultTimeEntry{
				Workspace: "w2",
				ProjectID: "p2",
				TaskID:    "t2",
				Billable:  &bTrue,
				TagIDs:    []string{"tg1", "tg2"},
			},
			expected: defaults.DefaultTimeEntry{
				Workspace: "w2",
				ProjectID: "p2",
				TaskID:    "t2",
				Billable:  &bFalse,
				TagIDs:    []string{"tg1", "tg2"},
			},
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)

			if len(tt.args) != 0 {
				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: false,
					Interactive:    false,
				})

				c := mocks.NewMockClient(t)

				tasks := make([]dto.Task, 0)

				if tt.expected.TaskID != "" {
					tasks = append(tasks, dto.Task{
						ID: tt.expected.TaskID,
					})
				}

				if tt.expected.ProjectID != "" {
					c.On("GetProject", mock.Anything).Return(&dto.Project{
						ID:    tt.expected.ProjectID,
						Tasks: tasks,
					}, nil)
				}

				if len(tt.expected.TagIDs) != 0 {
					tags := make([]dto.Tag, len(tt.expected.TagIDs))

					for i := range tt.expected.TagIDs {
						tags[i] = dto.Tag{ID: tt.expected.TagIDs[i]}
					}

					c.On("GetTags", mock.Anything).Return(tags, nil)
				}

				f.EXPECT().Client().Return(c, nil)
			}

			f.EXPECT().GetWorkspaceID().Return(tt.expected.Workspace, nil)

			ted := mocks.NewMockTimeEntryDefaults(t)
			ted.EXPECT().Read().Return(tt.current, nil)
			ted.EXPECT().Write(tt.expected).Return(nil)
			f.EXPECT().TimeEntryDefaults().Return(ted)

			result, called, err := runCmd(f, tt.args)

			assert.NoError(t, err, "should not have failed")
			assert.True(t, called)
			assert.Equal(t, tt.expected, result)
		})
	}
}
