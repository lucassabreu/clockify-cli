package util

import (
	"errors"
	"testing"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	. "github.com/lucassabreu/clockify-cli/internal/testhlp"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var bTrue = true
var bFalse = false

func TestDo_ShouldApplySteps_InOrder(t *testing.T) {
	first := func(TimeEntryDTO) (TimeEntryDTO, error) {
		return TimeEntryDTO{ID: "first"}, nil
	}
	second := func(TimeEntryDTO) (TimeEntryDTO, error) {
		return TimeEntryDTO{ID: "second"}, nil
	}

	tts := []struct {
		name   string
		steps  []Step
		result TimeEntryDTO
	}{
		{
			name: "first",
			steps: []Step{
				second,
				first,
				skip,
				skip,
			},
			result: TimeEntryDTO{ID: "first"},
		},
		{
			name: "second",
			steps: []Step{
				skip,
				first,
				skip,
				second,
			},
			result: TimeEntryDTO{ID: "second"},
		},
		{
			name: "only skips",
			steps: []Step{
				skip,
				skip,
				skip,
			},
			result: TimeEntryDTO{},
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			r, err := Do(TimeEntryDTO{}, tt.steps...)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.result, r)
		})
	}
}

func TestImplToDTOAndBack(t *testing.T) {
	end := MustParseTime(timehlp.SimplerTimeFormat, "2022-11-07 11:00")
	impl := dto.TimeEntryImpl{
		Billable:    true,
		Description: "unique",
		ID:          "te",
		IsLocked:    false,
		ProjectID:   "p",
		TagIDs:      []string{"tag"},
		TaskID:      "t",
		TimeInterval: dto.NewTimeInterval(
			MustParseTime(timehlp.SimplerTimeFormat, "2022-11-07 10:00"),
			&end,
		),
		UserID:      "u",
		WorkspaceID: "w",
	}

	dto := TimeEntryImplToDTO(impl)
	nimpl := TimeEntryDTOToImpl(dto)

	assert.Equal(t, impl, nimpl)
	assert.Equal(t, dto, TimeEntryImplToDTO(nimpl))
}

func TestTimeEntryDTOToImpl_ShouldFillMissingProperties(t *testing.T) {
	tm := MustParseTime(timehlp.SimplerTimeFormat, "2022-11-07 10:00")
	assert.Equal(t,
		dto.TimeEntryImpl{TimeInterval: dto.NewTimeInterval(tm, &tm)},
		TimeEntryDTOToImpl(TimeEntryDTO{
			Start: tm,
			End:   &tm,
		}),
	)
}

type flagSetMock struct {
	flags map[string]interface{}
}

func (f *flagSetMock) Changed(k string) bool {
	_, ok := f.flags[k]
	return ok
}

func (f *flagSetMock) GetString(k string) (string, error) {
	if f.Changed(k) {
		return f.flags[k].(string), nil
	}

	return "", nil
}

func (f *flagSetMock) GetStringSlice(k string) ([]string, error) {
	if f.Changed(k) {
		return f.flags[k].([]string), nil
	}

	return []string{}, nil
}

func TestFillTimeEntryWithFlags_ShouldNotSetProperties_WhenNotChanged(
	t *testing.T) {
	tm := MustParseTime(timehlp.SimplerTimeFormat, "2022-11-07 11:00").Local()

	forEnd := func(t time.Time) *time.Time { return &t }

	tts := []struct {
		name   string
		flags  flagSet
		input  TimeEntryDTO
		output TimeEntryDTO
		err    string
	}{
		{
			name: "only description",
			flags: &flagSetMock{flags: map[string]interface{}{
				"description": "diff",
			}},
			input:  TimeEntryDTO{Description: "other"},
			output: TimeEntryDTO{Description: "diff"},
		},
		{
			name: "only dates",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when":          tm.Format(timehlp.SimplerTimeFormat),
				"when-to-close": tm.Format(timehlp.SimplerTimeFormat),
			}},
			output: TimeEntryDTO{
				Start: tm,
				End:   &tm,
			},
		},
		{
			name: "should accept descriptive relative dates",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when":          "+1h15m15s",
				"when-to-close": "+2h16m16s",
			}},
			output: TimeEntryDTO{
				Start: timehlp.Now().
					Add(time.Hour + time.Minute*15 + time.Second*15),
				End: forEnd(
					timehlp.Now().
						Add(time.Hour*2 + time.Minute*16 + time.Second*16),
				),
			},
		},
		{
			name: "should accept relative dates",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when":          "+15:15",
				"when-to-close": "+16:16",
			}},
			output: TimeEntryDTO{
				Start: timehlp.Now().
					Add(time.Minute*15 + time.Second*15),
				End: forEnd(
					timehlp.Now().Add(time.Minute*16 + time.Second*16),
				),
			},
		},
		{
			name: "should accept fixed time",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when":          "15:15",
				"when-to-close": "16:16",
			}},
			output: TimeEntryDTO{
				Start: timehlp.Today().
					Add(time.Hour*15 + time.Minute*15),
				End: forEnd(
					timehlp.Today().Add(time.Hour*16 + time.Minute*16),
				),
			},
		},
		{
			name: "should validate time-strings",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when": "wrong",
			}},
			err: "parsing time.*",
		},
		{
			name: "should validate time-strings (close time)",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when-to-close": "wrong",
			}},
			err: "parsing time.*",
		},
		{
			name: "should validate time-strings with right format, but wrong",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when": "99:99:99",
			}},
			err: "parsing time.*",
		},
		{
			name: "should validate time-strings with right format, end time",
			flags: &flagSetMock{flags: map[string]interface{}{
				"when-to-close": "99:99:99",
			}},
			err: "parsing time.*",
		},
		{
			name: "all but dates and not-billable",
			flags: &flagSetMock{flags: map[string]interface{}{
				"description": "d",
				"project":     "p",
				"task":        "t",
				"tag":         []string{"t1", "t2"},
				"billable":    true,
			}},
			input: TimeEntryDTO{
				Start: tm,
				End:   &tm,
			},
			output: TimeEntryDTO{
				Start:       tm,
				End:         &tm,
				Description: "d",
				ProjectID:   "p",
				TaskID:      "t",
				TagIDs:      []string{"t1", "t2"},
				Billable:    &bTrue,
			},
		},
		{
			name: "all but dates and billable",
			flags: &flagSetMock{flags: map[string]interface{}{
				"description":  "d",
				"project":      "p",
				"task":         "t",
				"tags":         []string{"t1", "t2"},
				"not-billable": true,
			}},
			input: TimeEntryDTO{
				Start:    tm,
				End:      &tm,
				Billable: &bTrue,
			},
			output: TimeEntryDTO{
				Start:       tm,
				End:         &tm,
				Description: "d",
				ProjectID:   "p",
				TaskID:      "t",
				TagIDs:      []string{"t1", "t2"},
				Billable:    &bFalse,
			},
		},
		{
			name: "should reset task when project changes",
			flags: &flagSetMock{flags: map[string]interface{}{
				"description":  "d",
				"project":      "p2",
				"tags":         []string{"t1", "t2"},
				"not-billable": true,
			}},
			input: TimeEntryDTO{
				ProjectID: "p1",
				TaskID:    "t",
				Start:     tm,
				End:       &tm,
				Billable:  &bTrue,
			},
			output: TimeEntryDTO{
				Start:       tm,
				End:         &tm,
				Description: "d",
				ProjectID:   "p2",
				TaskID:      "",
				TagIDs:      []string{"t1", "t2"},
				Billable:    &bFalse,
			},
		},
		{
			name: "should not be billable and not billable",
			flags: &flagSetMock{flags: map[string]interface{}{
				"billable":     true,
				"not-billable": true,
			}},
			err: "flags can't be used together",
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			d, err := FillTimeEntryWithFlags(tt.flags)(tt.input)
			if tt.err != "" {
				if !assert.Error(t, err) {
					return
				}
				assert.Regexp(t, tt.err, err.Error())

				return
			}

			assert.Equal(t, tt.output, d)
		})
	}
}

func TestGetValidateTimeEntry_ShouldValidate_UsingSettingsAndConfigs(
	t *testing.T) {
	cnf := &mocks.SimpleConfig{
		AllowIncomplete: false,
	}

	wSettingsFn := func(
		ws dto.WorkspaceSettings) func(t *testing.T) cmdutil.Factory {
		return func(t *testing.T) cmdutil.Factory {
			f := mocks.NewMockFactory(t)
			f.EXPECT().Config().Return(cnf)

			f.EXPECT().GetWorkspace().Return(dto.Workspace{Settings: ws}, nil)

			return f
		}
	}

	wSettingsAndProjectFn := func(
		w dto.WorkspaceSettings,
		p *dto.Project,
		err error,
	) func(t *testing.T) cmdutil.Factory {
		return func(t *testing.T) cmdutil.Factory {
			f := wSettingsFn(w)(t).(*mocks.MockFactory)

			c := mocks.NewMockClient(t)
			f.EXPECT().Client().Return(c, nil)

			c.EXPECT().GetProject(mock.Anything).Return(p, err)

			return f
		}
	}

	tts := []struct {
		name    string
		input   TimeEntryDTO
		err     string
		factory func(*testing.T) cmdutil.Factory
	}{
		{
			name: "do nothing",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(&mocks.SimpleConfig{
					AllowIncomplete: true,
				})

				return f
			},
		},
		{
			name: "fail to find workspace",
			err:  "get workspace error",
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(cnf)

				f.EXPECT().GetWorkspace().
					Return(dto.Workspace{}, errors.New("get workspace error"))

				return f
			},
		},
		{
			name:  "fail to init client",
			err:   "client error",
			input: TimeEntryDTO{ProjectID: "p"},
			factory: func(t *testing.T) cmdutil.Factory {
				f := mocks.NewMockFactory(t)
				f.EXPECT().Config().Return(cnf)

				f.EXPECT().GetWorkspace().Return(dto.Workspace{}, nil)

				f.EXPECT().Client().
					Return(mocks.NewMockClient(t), errors.New("client error"))

				return f
			},
		},
		{
			name:  "force project",
			input: TimeEntryDTO{},
			err:   "workspace requires project",
			factory: wSettingsFn(dto.WorkspaceSettings{
				ForceProjects: true,
			}),
		},
		{
			name:  "force task",
			input: TimeEntryDTO{},
			err:   "workspace requires task",
			factory: wSettingsFn(dto.WorkspaceSettings{
				ForceTasks: true,
			}),
		},
		{
			name:  "force description",
			input: TimeEntryDTO{},
			err:   "workspace requires description",
			factory: wSettingsFn(dto.WorkspaceSettings{
				ForceDescription: true,
			}),
		},
		{
			name:  "force tags",
			input: TimeEntryDTO{},
			err:   "workspace requires at least one tag",
			factory: wSettingsFn(dto.WorkspaceSettings{
				ForceTags: true,
			}),
		},
		{
			name: "project not found",
			input: TimeEntryDTO{
				Workspace:   "w",
				Description: "description",
				ProjectID:   "project",
				TaskID:      "task",
				TagIDs:      []string{"tag"},
			},
			err: "not found",
			factory: wSettingsAndProjectFn(
				dto.WorkspaceSettings{
					ForceDescription: true,
					ForceProjects:    true,
					ForceTasks:       true,
					ForceTags:        true,
				},
				nil,
				api.EntityNotFound{
					EntityName: "project",
					ID:         "project",
				},
			),
		},
		{
			name: "project is archived, all is required",
			input: TimeEntryDTO{
				Workspace:   "w",
				Description: "description",
				ProjectID:   "project",
				TaskID:      "task",
				TagIDs:      []string{"tag"},
			},
			err: "project \\w+ - \\w+ is archived",
			factory: wSettingsAndProjectFn(
				dto.WorkspaceSettings{
					ForceDescription: true,
					ForceProjects:    true,
					ForceTasks:       true,
					ForceTags:        true,
				},
				&dto.Project{
					ID:       "project",
					Name:     "name",
					Archived: true,
				}, nil,
			),
		},
		{
			name: "project is archived, nothing is required",
			input: TimeEntryDTO{
				ProjectID: "project",
			},
			err: "project \\w+ - \\w+ is archived",
			factory: wSettingsAndProjectFn(
				dto.WorkspaceSettings{},
				&dto.Project{
					ID:       "project",
					Name:     "name",
					Archived: true,
				}, nil,
			),
		},
		{
			name: "nothing is required, without project",
			input: TimeEntryDTO{
				Workspace:   "w",
				Description: "description",
				TaskID:      "task",
				TagIDs:      []string{"tag"},
			},
			factory: wSettingsFn(dto.WorkspaceSettings{}),
		},
		{
			name: "everything is right",
			input: TimeEntryDTO{
				Workspace:   "w",
				Description: "description",
				ProjectID:   "project",
				TaskID:      "task",
				TagIDs:      []string{"tag"},
			},
			factory: wSettingsAndProjectFn(
				dto.WorkspaceSettings{
					ForceDescription: true,
					ForceProjects:    true,
					ForceTasks:       true,
					ForceTags:        true,
				},
				&dto.Project{
					ID:   "project",
					Name: "name",
				},
				nil,
			),
		},
		{
			name:    "nothing is required",
			input:   TimeEntryDTO{},
			factory: wSettingsFn(dto.WorkspaceSettings{}),
		},
	}

	for i := range tts {
		tt := &tts[i]
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetValidateTimeEntryFn(tt.factory(t))(tt.input)
			if tt.err != "" {
				if !assert.Error(t, err) {
					return
				}
				assert.Regexp(t, tt.err, err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetAllowNameForIDsFn_ShouldLookupEntityIDs_WhenFilled(t *testing.T) {
	te := TimeEntryDTO{
		Workspace: "w",
		ProjectID: "p",
		TaskID:    "t",
		TagIDs:    []string{"t1", "t2"},
	}

	cf := &mocks.SimpleConfig{AllowNameForID: true}
	c := mocks.NewMockClient(t)

	c.EXPECT().GetProjects(api.GetProjectsParam{
		Workspace:       te.Workspace,
		PaginationParam: api.AllPages(),
		Archived:        &bFalse,
	}).
		Return([]dto.Project{{ID: "pj_id", Name: "project"}}, nil)

	c.EXPECT().GetTasks(api.GetTasksParam{
		Workspace:       te.Workspace,
		ProjectID:       "pj_id",
		Active:          true,
		PaginationParam: api.AllPages(),
	}).
		Return([]dto.Task{{ID: "tk_id", Name: "task"}}, nil)

	c.EXPECT().GetTags(api.GetTagsParam{
		Workspace:       te.Workspace,
		Archived:        &bFalse,
		PaginationParam: api.AllPages(),
	}).
		Return([]dto.Tag{
			{ID: "tg_id_1", Name: "t1"},
			{ID: "tg_id_2", Name: "t2"},
		}, nil)

	te, err := GetAllowNameForIDsFn(cf, c)(te)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "pj_id", te.ProjectID)
	assert.Equal(t, "tk_id", te.TaskID)
	assert.Equal(t, []string{"tg_id_1", "tg_id_2"}, te.TagIDs)
}

func TestGetAllowNameForIDsFn_ShouldNotLookupEntityIDs_WhenEmpty(
	t *testing.T) {
	te := TimeEntryDTO{
		Workspace: "",
		ProjectID: "",
		TaskID:    "",
		TagIDs:    []string{},
	}

	te2, err := GetAllowNameForIDsFn(
		&mocks.SimpleConfig{AllowNameForID: true}, mocks.NewMockClient(t))(te)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, te, te2)
}

func TestGetAllowNameForIDsFn_ShouldNotLookup_WhenDisabled(t *testing.T) {
	te := TimeEntryDTO{
		Workspace: "w",
		ProjectID: "p",
		TaskID:    "t",
		TagIDs:    []string{"t1", "t2"},
	}

	te2, err := GetAllowNameForIDsFn(
		&mocks.SimpleConfig{AllowNameForID: false}, mocks.NewMockClient(t))(te)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, te, te2)
}

func TestGetAllowNameForIDsFn_ShouldFail_WhenEntitiesNotFound(t *testing.T) {
	te := TimeEntryDTO{
		Workspace: "w",
		ProjectID: "p",
		TaskID:    "t",
		TagIDs:    []string{"t1"},
	}

	cf := &mocks.SimpleConfig{AllowNameForID: true}
	c := mocks.NewMockClient(t)

	c.EXPECT().GetProjects(mock.Anything).Return([]dto.Project{}, nil).Once()
	c.EXPECT().GetTasks(mock.Anything).Return([]dto.Task{}, nil).Once()
	c.EXPECT().GetTags(mock.Anything).Return([]dto.Tag{}, nil).Once()

	_, err := GetAllowNameForIDsFn(cf, c)(te)
	if !assert.Error(t, err) {
		return
	}
	assert.Regexp(t, "No project with id or name .*", err.Error())

	c.EXPECT().GetProjects(api.GetProjectsParam{
		Workspace:       te.Workspace,
		Archived:        &bFalse,
		PaginationParam: api.AllPages(),
	}).
		Return([]dto.Project{{ID: "pj_id", Name: "project"}}, nil)

	_, err = GetAllowNameForIDsFn(cf, c)(te)
	if !assert.Error(t, err) {
		return
	}
	assert.Regexp(t, "No task with id or name .*", err.Error())

	c.EXPECT().GetTasks(api.GetTasksParam{
		Workspace:       te.Workspace,
		ProjectID:       "pj_id",
		Active:          true,
		PaginationParam: api.AllPages(),
	}).
		Return([]dto.Task{{ID: "tk_id", Name: "task"}}, nil)

	_, err = GetAllowNameForIDsFn(cf, c)(te)
	if !assert.Error(t, err) {
		return
	}
	assert.Regexp(t, "No tag with id or name .*", err.Error())

	c.EXPECT().GetTags(api.GetTagsParam{
		Workspace:       te.Workspace,
		Archived:        &bFalse,
		PaginationParam: api.AllPages(),
	}).
		Return([]dto.Tag{{ID: "tg_id_1", Name: "t1"}}, nil)

	_, err = GetAllowNameForIDsFn(cf, c)(te)
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, err)
}

func TestGetAllowNameForIDsFn_ShouldBeQuiet_WhenInteractive(t *testing.T) {
	te := TimeEntryDTO{
		Workspace: "w",
		ProjectID: "p",
		TagIDs:    []string{"t1"},
	}

	cf := &mocks.SimpleConfig{
		AllowNameForID: true,
		Interactive:    true,
	}
	c := mocks.NewMockClient(t)

	c.EXPECT().GetProjects(mock.Anything).Return([]dto.Project{}, nil).Once()
	c.EXPECT().GetTags(mock.Anything).Return([]dto.Tag{}, nil).Once()

	_, err := GetAllowNameForIDsFn(cf, c)(te)
	if !assert.NoError(t, err) {
		return
	}

	te.TaskID = "t"

	c.EXPECT().GetProjects(mock.Anything).
		Return([]dto.Project{{ID: "pj_id", Name: "project"}}, nil).Once()

	c.EXPECT().GetTasks(mock.Anything).Return([]dto.Task{}, nil).Once()
	c.EXPECT().GetTags(mock.Anything).Return([]dto.Tag{}, nil).Once()

	_, err = GetAllowNameForIDsFn(cf, c)(te)
	if !assert.NoError(t, err) {
		return
	}
}
