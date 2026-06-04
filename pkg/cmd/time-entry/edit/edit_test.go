package edit_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/edit"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdEditWhenChangingProjectOrTask(t *testing.T) {
	w := dto.Workspace{ID: "w"}
	te := dto.TimeEntryImpl{
		WorkspaceID: w.ID,
		ID:          "timeentryid",
		Description: "Something",
		ProjectID:   "oldproj",
		TaskID:      "oldtask",
		TimeInterval: dto.TimeInterval{
			Start: time.Now(),
		},
	}

	tts := []struct {
		name        string
		args        []string
		project     *dto.Project
		updateParam api.UpdateTimeEntryParam
	}{
		{
			name:    "should remove task, when changing project",
			args:    []string{"-p", "newproj"},
			project: &dto.Project{ID: "newproj", Name: "newproj"},
			updateParam: api.UpdateTimeEntryParam{
				Workspace:   te.WorkspaceID,
				TimeEntryID: te.ID,
				Start:       te.TimeInterval.Start,
				End:         te.TimeInterval.End,
				Billable:    te.Billable,
				Description: te.Description,
				ProjectID:   "newproj",
				TaskID:      "",
				TagIDs:      te.TagIDs,
			},
		},
		{
			name: "should remove task, when removing project",
			args: []string{"-p", ""},
			updateParam: api.UpdateTimeEntryParam{
				Workspace:   te.WorkspaceID,
				TimeEntryID: te.ID,
				Start:       te.TimeInterval.Start,
				End:         te.TimeInterval.End,
				Billable:    te.Billable,
				Description: te.Description,
				ProjectID:   "",
				TaskID:      "",
				TagIDs:      te.TagIDs,
			},
		},
		{
			name:    "should change project and task",
			args:    []string{"--task", "newtask", "-p=newproj"},
			project: &dto.Project{ID: "newproj", Name: "newproj"},
			updateParam: api.UpdateTimeEntryParam{
				Workspace:   te.WorkspaceID,
				TimeEntryID: te.ID,
				Start:       te.TimeInterval.Start,
				End:         te.TimeInterval.End,
				Billable:    te.Billable,
				Description: te.Description,
				ProjectID:   "newproj",
				TaskID:      "newtask",
				TagIDs:      te.TagIDs,
			},
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)

			f.EXPECT().GetUserID().Return("u", nil)
			f.EXPECT().GetWorkspace().Return(w, nil)
			f.EXPECT().GetWorkspaceID().Return(w.ID, nil)

			f.EXPECT().Config().Return(&mocks.SimpleConfig{
				AllowNameForID: true,
			})

			c := mocks.NewMockClient(t)
			f.EXPECT().Client().Return(c, nil)

			c.EXPECT().GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
				Workspace: "w",
				UserID:    "u",
			}).
				Return(&te, nil)

			p := tt.project
			if p != nil {
				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       w.ID,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{*p}, nil)

				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: w.ID,
					ProjectID: p.ID,
				}).
					Return(p, nil)
			}

			if tt.updateParam.TaskID != "" {
				c.EXPECT().GetTasks(api.GetTasksParam{
					Workspace:       w.ID,
					ProjectID:       tt.updateParam.ProjectID,
					Active:          true,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Task{{ID: tt.updateParam.TaskID}}, nil)
			}

			c.EXPECT().UpdateTimeEntry(tt.updateParam).
				Return(te, nil)

			called := false
			cmd := edit.NewCmdEdit(f, func(
				_ dto.TimeEntryImpl, _ io.Writer, _ util.OutputFlags) error {
				called = true
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(append(tt.args, "current", "-q"))
			_, err := cmd.ExecuteC()

			if assert.NoError(t, err) {
				t.Cleanup(func() {
					assert.True(t, called)
				})
				return
			}

			t.Fatalf("err: %s", err)
		})
	}
}

func TestNewCmdEditMultipleTimeEntries(t *testing.T) {
	w := dto.Workspace{ID: "w"}
	now := time.Now()
	start1 := now
	end1 := now.Add(1 * time.Hour)
	start2 := now.Add(2 * time.Hour)
	end2 := now.Add(3 * time.Hour)

	te1 := dto.TimeEntryImpl{
		WorkspaceID: w.ID,
		ID:          "teid1",
		Description: "Entry 1",
		ProjectID:   "proj1",
		TaskID:      "task1",
		TagIDs:      []string{"tag1"},
		TimeInterval: dto.TimeInterval{
			Start: start1,
			End:   &end1,
		},
		Billable: false,
		UserID:   "u",
	}

	te2 := dto.TimeEntryImpl{
		WorkspaceID: w.ID,
		ID:          "teid2",
		Description: "Entry 2",
		ProjectID:   "proj1",
		TaskID:      "task2",
		TagIDs:      []string{"tag2"},
		TimeInterval: dto.TimeInterval{
			Start: start2,
			End:   &end2,
		},
		Billable: true,
		UserID:   "u",
	}

	tts := []struct {
		name            string
		args            []string
		timeEntries     []dto.TimeEntryImpl
		updateParams    []api.UpdateTimeEntryParam
		validateProject *dto.Project
		err             string
	}{
		{
			name: "should fail when using --when with multiple entries",
			args: []string{"teid1", "teid2", "-s", "09:00"},
			err:  "--when and --when-to-close can only be used when editing a single time entry",
		},
		{
			name: "should fail when using --when-to-close with multiple entries",
			args: []string{"teid1", "teid2", "-e", "18:00"},
			err:  "--when and --when-to-close can only be used when editing a single time entry",
		},
		{
			name:        "should update only description",
			args:        []string{"teid1", "teid2", "-d", "New Description", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1, te2},
			updateParams: []api.UpdateTimeEntryParam{
				{
					Workspace:   te1.WorkspaceID,
					TimeEntryID: te1.ID,
					Start:       te1.TimeInterval.Start,
					End:         te1.TimeInterval.End,
					Billable:    te1.Billable,
					Description: "New Description",
					ProjectID:   te1.ProjectID,
					TaskID:      te1.TaskID,
					TagIDs:      te1.TagIDs,
				},
				{
					Workspace:   te2.WorkspaceID,
					TimeEntryID: te2.ID,
					Start:       te2.TimeInterval.Start,
					End:         te2.TimeInterval.End,
					Billable:    te2.Billable,
					Description: "New Description",
					ProjectID:   te2.ProjectID,
					TaskID:      te2.TaskID,
					TagIDs:      te2.TagIDs,
				},
			},
			validateProject: &dto.Project{ID: te1.ProjectID},
		},
		{
			name:        "should update only tags",
			args:        []string{"teid1", "teid2", "-T", "newtag", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1, te2},
			updateParams: []api.UpdateTimeEntryParam{
				{
					Workspace:   te1.WorkspaceID,
					TimeEntryID: te1.ID,
					Start:       te1.TimeInterval.Start,
					End:         te1.TimeInterval.End,
					Billable:    te1.Billable,
					Description: te1.Description,
					ProjectID:   te1.ProjectID,
					TaskID:      te1.TaskID,
					TagIDs:      []string{"newtag"},
				},
				{
					Workspace:   te2.WorkspaceID,
					TimeEntryID: te2.ID,
					Start:       te2.TimeInterval.Start,
					End:         te2.TimeInterval.End,
					Billable:    te2.Billable,
					Description: te2.Description,
					ProjectID:   te2.ProjectID,
					TaskID:      te2.TaskID,
					TagIDs:      []string{"newtag"},
				},
			},
			validateProject: &dto.Project{ID: te1.ProjectID},
		},
		{
			name:        "should update only task",
			args:        []string{"teid1", "teid2", "--task", "newtask", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1, te2},
			updateParams: []api.UpdateTimeEntryParam{
				{
					Workspace:   te1.WorkspaceID,
					TimeEntryID: te1.ID,
					Start:       te1.TimeInterval.Start,
					End:         te1.TimeInterval.End,
					Billable:    te1.Billable,
					Description: te1.Description,
					ProjectID:   te1.ProjectID,
					TaskID:      "newtask",
					TagIDs:      te1.TagIDs,
				},
				{
					Workspace:   te2.WorkspaceID,
					TimeEntryID: te2.ID,
					Start:       te2.TimeInterval.Start,
					End:         te2.TimeInterval.End,
					Billable:    te2.Billable,
					Description: te2.Description,
					ProjectID:   te1.ProjectID,
					TaskID:      "newtask",
					TagIDs:      te2.TagIDs,
				},
			},
			validateProject: &dto.Project{ID: te1.ProjectID},
		},
		{
			name:        "should clear task when changing project",
			args:        []string{"teid1", "teid2", "-p", "newproj", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1, te2},
			updateParams: []api.UpdateTimeEntryParam{
				{
					Workspace:   te1.WorkspaceID,
					TimeEntryID: te1.ID,
					Start:       te1.TimeInterval.Start,
					End:         te1.TimeInterval.End,
					Billable:    te1.Billable,
					Description: te1.Description,
					ProjectID:   "newproj",
					TaskID:      "",
					TagIDs:      te1.TagIDs,
				},
				{
					Workspace:   te2.WorkspaceID,
					TimeEntryID: te2.ID,
					Start:       te2.TimeInterval.Start,
					End:         te2.TimeInterval.End,
					Billable:    te2.Billable,
					Description: te2.Description,
					ProjectID:   "newproj",
					TaskID:      "",
					TagIDs:      te2.TagIDs,
				},
			},
			validateProject: &dto.Project{ID: "newproj"},
		},
		{
			name:        "should clear both project and task when setting project to empty",
			args:        []string{"teid1", "teid2", "-p", "", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1, te2},
			updateParams: []api.UpdateTimeEntryParam{
				{
					Workspace:   te1.WorkspaceID,
					TimeEntryID: te1.ID,
					Start:       te1.TimeInterval.Start,
					End:         te1.TimeInterval.End,
					Billable:    te1.Billable,
					Description: te1.Description,
					ProjectID:   "",
					TaskID:      "",
					TagIDs:      te1.TagIDs,
				},
				{
					Workspace:   te2.WorkspaceID,
					TimeEntryID: te2.ID,
					Start:       te2.TimeInterval.Start,
					End:         te2.TimeInterval.End,
					Billable:    te2.Billable,
					Description: te2.Description,
					ProjectID:   "",
					TaskID:      "",
					TagIDs:      te2.TagIDs,
				},
			},
		},
		{
			name: "should fail when changing task on entries with different projects",
			args: []string{"teid1", "teid2", "--task", "newtask", "-q"},
			timeEntries: []dto.TimeEntryImpl{te1,

				func() dto.TimeEntryImpl {
					t := te2
					t.ProjectID = "proj2"
					return t
				}(),
			},
			validateProject: &dto.Project{ID: te1.ProjectID},
			err:             "you are changing the task of the time entries, but not the project and some of them are not in the same project, please also set --project",
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)

			f.EXPECT().GetUserID().Return("u", nil)
			f.EXPECT().GetWorkspaceID().Return(w.ID, nil)

			c := mocks.NewMockClient(t)
			f.EXPECT().Client().Return(c, nil)

			f.EXPECT().Config().Return(&mocks.SimpleConfig{AllowNameForID: false})

			if len(tt.timeEntries) > 0 {
				f.EXPECT().GetWorkspace().Return(w, nil)

				for i, te := range tt.timeEntries {
					c.EXPECT().GetTimeEntry(api.GetTimeEntryParam{
						Workspace:   w.ID,
						TimeEntryID: te.ID,
					}).Return(&tt.timeEntries[i], nil)
				}
			}

			if tt.validateProject != nil {
				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: w.ID,
					ProjectID: tt.validateProject.ID,
				}).Return(tt.validateProject, nil)
			}

			if tt.err == "" {
				for i := range tt.timeEntries {
					tei := tt.timeEntries[i]
					c.EXPECT().GetHydratedTimeEntry(api.GetTimeEntryParam{
						Workspace:   w.ID,
						TimeEntryID: tei.ID,
					}).Return(&dto.TimeEntry{ID: tei.ID, WorkspaceID: w.ID}, nil)
				}

				for _, p := range tt.updateParams {
					c.EXPECT().UpdateTimeEntry(p).Return(dto.TimeEntryImpl{}, nil)
				}
			}

			cmd := edit.NewCmdEdit(f, func(
				_ dto.TimeEntryImpl, _ io.Writer, _ util.OutputFlags) error {
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if tt.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
