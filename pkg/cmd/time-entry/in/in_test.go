package in_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/in"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
)

var w = dto.Workspace{ID: "w"}

func TestNewCmdIn_ShouldBeBothBillableAndNotBillable(t *testing.T) {
	f := mocks.NewMockFactory(t)

	f.EXPECT().GetUserID().Return("u", nil)
	f.EXPECT().GetWorkspaceID().Return(w.ID, nil)

	f.EXPECT().Config().Return(&mocks.SimpleConfig{})

	c := mocks.NewMockClient(t)
	f.EXPECT().Client().Return(c, nil)

	called := false
	cmd := in.NewCmdIn(f, func(
		_ dto.TimeEntryImpl, _ io.Writer, _ util.OutputFlags) error {
		called = true
		return nil
	})

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	out := bytes.NewBufferString("")
	cmd.SetOut(out)
	cmd.SetErr(out)

	cmd.SetArgs([]string{"--billable", "--not-billable"})
	_, err := cmd.ExecuteC()

	if assert.Error(t, err) {
		assert.False(t, called)
		flagErr := &cmdutil.FlagError{}
		assert.ErrorAs(t, err, &flagErr)
		return
	}

	t.Fatal("should've failed")
}

func TestNewCmdIn_ShouldNotSetBillable_WhenNotAsked(t *testing.T) {
	bTrue := true
	bFalse := false

	tts := []struct {
		name  string
		args  []string
		param api.CreateTimeEntryParam
	}{
		{
			name: "should be nil",
			args: []string{"-s=08:00"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     timehlp.Today().Add(8 * time.Hour),
				Billable:  nil,
			},
		},
		{
			name: "should be billable",
			args: []string{"-s=08:00", "--billable"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     timehlp.Today().Add(8 * time.Hour),
				Billable:  &bTrue,
			},
		},
		{
			name: "should not be billable",
			args: []string{"-s=08:00", "--not-billable"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     timehlp.Today().Add(8 * time.Hour),
				Billable:  &bFalse,
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
				Workspace: w.ID,
				UserID:    "u",
			}).
				Return(nil, nil)

			c.EXPECT().Out(api.OutParam{
				Workspace: w.ID,
				UserID:    "u",
				End:       tt.param.Start,
			}).Return(api.ErrorNotFound)

			c.EXPECT().CreateTimeEntry(tt.param).
				Return(dto.TimeEntryImpl{ID: "te"}, nil)

			called := false
			cmd := in.NewCmdIn(f, func(
				_ dto.TimeEntryImpl, _ io.Writer, _ util.OutputFlags) error {
				called = true
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(append(tt.args, "-q"))
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

func TestNewCmdIn_ShouldLookupProject_WithAndWithoutClient(t *testing.T) {
	defaultStart := timehlp.Today().Add(8 * time.Hour)

	projects := []dto.Project{
		{ID: "p1", Name: "first", ClientID: "c1", ClientName: "other"},
		{ID: "p2", Name: "second", ClientID: "c2", ClientName: "me"},
		{ID: "p3", Name: "second", ClientID: "c3", ClientName: "clockify"},
		{ID: "p4", Name: "third"},
		{ID: "p5", Name: "notonclient", ClientID: "c3", ClientName: "clockify"},
	}

	tts := []struct {
		name  string
		args  []string
		param api.CreateTimeEntryParam
		err   error
	}{
		{
			name: "only project",
			args: []string{"-s=08:00", "-p=first"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[0].ID,
			},
		},
		{
			name: "project and client",
			args: []string{"-s=08:00", "-p=second", "-c=me"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[1].ID,
			},
		},
		{
			name: "project and other client",
			args: []string{"-s=08:00", "-p=second", "-c=clockify"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[2].ID,
			},
		},
		{
			name: "project without client",
			args: []string{"-s=08:00", "-p=third"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[3].ID,
			},
		},
		{
			name: "project does not exist",
			args: []string{"-s=08:00", "-p=notfound"},
			err: errors.New(
				"No project with id or name containing 'notfound' " +
					"was found"),
		},
		{
			name: "project does not exist in this client",
			args: []string{"-s=08:00", "-p=notonclient", "-c=me"},
			err: errors.New(
				"No project with id or name containing 'notonclient' " +
					"was found for client 'me'"),
		},
		{
			name: "project with client name does not exist",
			args: []string{"-s=08:00", "-p", "notonclient me"},
			err: errors.New(
				"No project with id or name containing 'notonclient me' " +
					"was found"),
		},
		{
			name: "project and client's name",
			args: []string{"-s=08:00", "-p", "second me"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[1].ID,
			},
		},
		{
			name: "project and client's name (other)",
			args: []string{"-s=08:00", "-p=second clockify"},
			param: api.CreateTimeEntryParam{
				Workspace: w.ID,
				Start:     defaultStart,
				ProjectID: projects[2].ID,
			},
		},
	}

	for i := range tts {
		tt := &tts[i]

		t.Run(tt.name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)

			f.EXPECT().GetUserID().Return("u", nil)
			f.EXPECT().GetWorkspaceID().Return(w.ID, nil)

			f.EXPECT().Config().Return(&mocks.SimpleConfig{
				AllowNameForID:               true,
				SearchProjectWithClientsName: true,
			})

			c := mocks.NewMockClient(t)
			f.EXPECT().Client().Return(c, nil)

			c.EXPECT().GetProjects(api.GetProjectsParam{
				Workspace:       w.ID,
				PaginationParam: api.AllPages(),
			}).
				Return(projects, nil)

			c.EXPECT().GetTimeEntryInProgress(api.GetTimeEntryInProgressParam{
				Workspace: w.ID,
				UserID:    "u",
			}).
				Return(nil, nil)

			if tt.err == nil {
				c.EXPECT().GetProject(api.GetProjectParam{
					Workspace: w.ID,
					ProjectID: tt.param.ProjectID,
				}).
					Return(&dto.Project{ID: tt.param.ProjectID}, nil)

				f.EXPECT().GetWorkspace().Return(w, nil)

				c.EXPECT().Out(api.OutParam{
					Workspace: w.ID,
					UserID:    "u",
					End:       tt.param.Start,
				}).Return(api.ErrorNotFound)

				c.EXPECT().CreateTimeEntry(tt.param).
					Return(dto.TimeEntryImpl{ID: "te"}, nil)
			}

			called := false
			cmd := in.NewCmdIn(f, func(
				_ dto.TimeEntryImpl, _ io.Writer, _ util.OutputFlags) error {
				called = true
				return nil
			})

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			out := bytes.NewBufferString("")
			cmd.SetOut(out)
			cmd.SetErr(out)

			cmd.SetArgs(append(tt.args, "-q"))
			_, err := cmd.ExecuteC()

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
				return
			}

			t.Cleanup(func() {
				assert.True(t, called)
			})

			if assert.NoError(t, err) {
				return
			}

			t.Fatalf("err: %s", err)
		})
	}

}
