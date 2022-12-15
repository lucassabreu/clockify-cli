package init_test

import (
	"strings"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/internal/consoletest"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	ini "github.com/lucassabreu/clockify-cli/pkg/cmd/config/init"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setStringFn(config *mocks.MockConfig, name, value string) *mock.Call {
	r := ""
	config.On("GetString", name).
		Return(func(string) string {
			v := r
			r = value
			return v
		})
	return config.On("SetString", name, value)
}

func setBoolFn(config *mocks.MockConfig, name string, first, value bool) *mock.Call {
	r := first
	config.On("GetBool", name).
		Return(func(string) bool {
			v := r
			r = value
			return v
		})
	return config.On("SetBool", name, value)
}

func TestInitCmd(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui.SetDefaultOptions(survey.WithStdio(in, out, out))

			f := mocks.NewMockFactory(t)
			config := mocks.NewMockConfig(t)
			client := mocks.NewMockClient(t)

			f.On("Config").Return(config)

			f.On("Client").NotBefore(
				config.On("GetString", cmdutil.CONF_TOKEN).Return(""),
				config.On("SetString", cmdutil.CONF_TOKEN, "new token"),
			).Return(client, nil)

			client.On("GetWorkspaces", api.GetWorkspaces{}).
				Return([]dto.Workspace{
					{ID: "1", Name: "First"},
					{ID: "2", Name: "Second"},
				}, nil)

			call := setStringFn(config, cmdutil.CONF_WORKSPACE, "2")

			client.On("WorkspaceUsers", api.WorkspaceUsersParam{
				Workspace:       "2",
				PaginationParam: api.AllPages(),
			}).
				NotBefore(call).
				Return([]dto.User{
					{ID: "user-1", Name: "John Due"},
					{ID: "user-2", Name: "Joana D'ark"},
				}, nil)

			setStringFn(config, cmdutil.CONF_USER_ID, "user-1")

			setBoolFn(config, cmdutil.CONF_ALLOW_NAME_FOR_ID, false, true)
			setBoolFn(config, cmdutil.CONF_INTERACTIVE, false, false)

			config.EXPECT().GetInt(cmdutil.CONF_INTERACTIVE_PAGE_SIZE).
				Return(7)
			config.EXPECT().
				SetInt(cmdutil.CONF_INTERACTIVE_PAGE_SIZE, 10)

			config.On("GetStringSlice", cmdutil.CONF_WORKWEEK_DAYS).
				Return([]string{})
			config.On("SetStringSlice", cmdutil.CONF_WORKWEEK_DAYS, []string{
				strings.ToLower(time.Sunday.String()),
				strings.ToLower(time.Tuesday.String()),
				strings.ToLower(time.Thursday.String()),
				strings.ToLower(time.Friday.String()),
				strings.ToLower(time.Saturday.String()),
			})

			setBoolFn(config, cmdutil.CONF_ALLOW_INCOMPLETE, false, false)
			setBoolFn(config, cmdutil.CONF_SHOW_TASKS, true, true)
			setBoolFn(config, cmdutil.CONF_SHOW_TOTAL_DURATION, true, true)
			setBoolFn(config, cmdutil.CONF_DESCR_AUTOCOMP, false, true)

			config.On("GetInt", cmdutil.CONF_DESCR_AUTOCOMP_DAYS).Return(0)
			config.On("SetInt", cmdutil.CONF_DESCR_AUTOCOMP_DAYS, 10)

			setBoolFn(config, cmdutil.CONF_ALLOW_ARCHIVED_TAGS, true, false)

			config.On("Save").Once().Return(nil)

			cmd := ini.NewCmdInit(f)
			_, err := cmd.ExecuteC()
			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("Token: ")
			c.SendLine("new token")
			c.ExpectString("new token")

			c.ExpectString("Choose default Workspace:")
			c.ExpectString("First")
			c.ExpectString("Second")
			c.SendLine("sec")
			c.ExpectString("Second")

			c.ExpectString("Choose your user:")
			c.ExpectString("John Due")
			c.ExpectString("Joana")
			c.SendLine("due")
			c.ExpectString("John Due")

			c.ExpectString("Should try to find")
			c.ExpectString("by their names?")
			c.SendLine("y")
			c.ExpectString("Yes")

			c.ExpectString("Interactive Mode\" by default?")
			c.SendLine("n")
			c.ExpectString("No")

			c.ExpectString("How many items should be shown when asking for " +
				"projects, tasks or tags?")
			c.ExpectString("7")
			c.SendLine("10")

			c.ExpectString("Which days of the week do you work?")
			c.ExpectString("sunday")
			c.ExpectString("monday")
			c.ExpectString("tuesday")
			c.ExpectString("wednesday")
			c.ExpectString("thursday")
			c.ExpectString("friday")
			c.ExpectString("saturday")

			c.Send(string(terminal.KeySpace))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeySpace))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeySpace))
			c.Send(string(terminal.KeyArrowUp))
			c.Send(string(terminal.KeySpace))
			c.Send("sat")
			c.Send(string(terminal.KeySpace))
			c.SendLine("")
			c.ExpectString("sunday, tuesday, thursday, friday, saturday")

			c.ExpectString("incomplete data?")
			c.SendLine("")
			c.ExpectString("No")

			c.ExpectString("show task on time entries")
			c.SendLine("")
			c.ExpectString("Yes")

			c.ExpectString("sum of the time entries duration?")
			c.SendLine("yes")
			c.ExpectString("Yes")

			c.ExpectString("descriptions?")
			c.SendLine("YES")
			c.ExpectString("Yes")

			c.ExpectString("How many days")
			c.SendLine("10")

			c.ExpectString("archived tags?")
			c.SendLine("n")
			c.ExpectString("No")

			c.ExpectEOF()
		})
}
func TestInitCmdCtrlC(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			ui.SetDefaultOptions(survey.WithStdio(in, out, out))

			f := mocks.NewMockFactory(t)
			config := mocks.NewMockConfig(t)

			f.On("Config").Return(config)
			config.On("GetString", cmdutil.CONF_TOKEN).Return("")

			cmd := ini.NewCmdInit(f)
			_, err := cmd.ExecuteC()
			if !assert.Error(t, err) {
				return nil
			}

			assert.ErrorIs(t, err, terminal.InterruptErr)
			return nil
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("Token: ")
			c.SendLine(string(terminal.KeyInterrupt))

			c.ExpectEOF()
		})
}
