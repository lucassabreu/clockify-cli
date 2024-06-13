package cmdcomplutil_test

import (
	"errors"
	"testing"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewProjectAutoComplete(t *testing.T) {
	bFalse := false
	tts := []struct {
		name       string
		toComplete string
		factory    func(t *testing.T) cmdutil.Factory
		err        string
		args       cmdcompl.ValidArgs
	}{
		{
			name: "no workspace, nothing",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().
					Return("", errors.New("no workspace"))
				return f
			},
			err:  "no workspace",
			args: cmdcompl.EmptyValidArgs(),
		},
		{
			name: "no client, nothing",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Client().Return(nil, errors.New("no client"))
				return f
			},
			err:  "no client",
			args: cmdcompl.EmptyValidArgs(),
		},
		{
			name: "fail get projects, nothing",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(mock.Anything).
					Return(nil, errors.New("fail to request"))

				return f
			},
			err:  "fail to request",
			args: cmdcompl.EmptyValidArgs(),
		},
		{
			name: "all projects",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{
						{ID: "p1", Name: "Project 1"},
						{ID: "p2", Name: "Project 2"},
						{ID: "p3", Name: "Project 3"},
					}, nil)

				return f
			},
			args: cmdcompl.ValidArgsMap{
				"p1": "Project 1",
				"p2": "Project 2",
				"p3": "Project 3",
			},
		},
		{
			name: "only projects with id with cat",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{})
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{
						{ID: "p0dog", Name: "Project 0"},
						{ID: "pcat", Name: "Project 1"},
						{ID: "catp", Name: "Project 2"},
						{ID: "pcatp", Name: "Project 3"},
						{ID: "p4", Name: "Project 4"},
					}, nil)

				return f
			},
			toComplete: "cat",
			args: cmdcompl.ValidArgsMap{
				"pcat":  "Project 1",
				"catp":  "Project 2",
				"pcatp": "Project 3",
			},
		},
		{
			name: "only projects with id or name with cat",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID: true,
				})

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{
						{ID: "p0dog", Name: "Project 0"},
						{ID: "pcat", Name: "Project 1"},
						{ID: "catp", Name: "Project 2"},
						{ID: "pcatp", Name: "Project 3"},
						{ID: "p4", Name: "Project 4"},
						{ID: "p5", Name: "Project Cat 5"},
						{ID: "p6", Name: "CAT Project 6"},
						{ID: "p7", Name: "Project 7 cats"},
						{ID: "p8", Name: "Catalog"},
					}, nil)

				return f
			},
			toComplete: "cat",
			args: cmdcompl.ValidArgsMap{
				"pcat":  "Project 1",
				"catp":  "Project 2",
				"pcatp": "Project 3",
				"p5":    "Project Cat 5",
				"p6":    "CAT Project 6",
				"p7":    "Project 7 cats",
				"p8":    "Catalog",
			},
		},
		{
			name: "only projects where the client has id or name with cat",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().GetWorkspaceID().Return("w", nil)

				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowNameForID:               true,
					SearchProjectWithClientsName: true,
				})

				c := mocks.NewMockClient(t)
				f.EXPECT().Client().Return(c, nil)

				c.EXPECT().GetProjects(api.GetProjectsParam{
					Workspace:       "w",
					Archived:        &bFalse,
					PaginationParam: api.AllPages(),
				}).
					Return([]dto.Project{
						{ID: "p1", Name: "Project 1"},
						{ID: "p2", Name: "CAT Project"},
						{ID: "p3", Name: "Catalog"},
						{ID: "p4", Name: "Project",
							ClientID: "c10", ClientName: "Cats"},
						{ID: "p5", Name: "Project",
							ClientID: "cat", ClientName: "Client"},
						{ID: "p6", Name: "Project",
							ClientID: "c30", ClientName: "Client"},
					}, nil)

				return f
			},
			toComplete: "cat",
			args: cmdcompl.ValidArgsMap{
				"p2": "CAT Project | Without Client",
				"p3": "Catalog     | Without Client",
				"p4": "Project     | c10 -- Cats",
				"p5": "Project     | cat -- Client",
			},
		},
	}

	for i := range tts {
		tt := tts[i]
		t.Run(tt.name, func(t *testing.T) {
			f := tt.factory(t)
			autoComplete := cmdcomplutil.NewProjectAutoComplete(
				f, f.Config())

			args, err := autoComplete(
				&cobra.Command{}, []string{}, tt.toComplete)

			if tt.err == "" && !assert.NoError(t, err) {
				return
			}

			if tt.err != "" && (!assert.Error(t, err) || !assert.Regexp(
				t, tt.err, err.Error())) {
				return
			}

			assert.Equal(t, tt.args, args)
		})
	}

}
