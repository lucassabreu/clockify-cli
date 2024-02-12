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
				bFalse := false
				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       w.ID,
					Archived:        &bFalse,
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
