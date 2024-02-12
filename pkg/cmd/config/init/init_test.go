package init_test

import (
	"errors"
	"strings"
	"testing"
	"time"

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
	return setStringDefaultFn(config, name, "", value)
}

func setStringDefaultFn(
	config *mocks.MockConfig, name, first, value string) *mock.Call {
	r := first
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
			f := mocks.NewMockFactory(t)
			config := mocks.NewMockConfig(t)
			client := mocks.NewMockClient(t)

			f.EXPECT().Config().Return(config)
			config.EXPECT().GetString(cmdutil.CONF_TOKEN).Return("")
			config.EXPECT().SetString(cmdutil.CONF_TOKEN, "new token")

			f.EXPECT().Client().Return(client, nil)

			client.EXPECT().GetWorkspaces(api.GetWorkspaces{}).
				Return([]dto.Workspace{
					{ID: "1", Name: "First"},
					{ID: "2", Name: "Second"},
				}, nil)

			call := setStringFn(config, cmdutil.CONF_WORKSPACE, "2")

			client.EXPECT().WorkspaceUsers(api.WorkspaceUsersParam{
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

			config.EXPECT().GetStringSlice(cmdutil.CONF_WORKWEEK_DAYS).
				Return([]string{})
			config.EXPECT().SetStringSlice(cmdutil.CONF_WORKWEEK_DAYS, []string{
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

			config.EXPECT().GetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS).Return(0)
			config.EXPECT().SetInt(cmdutil.CONF_DESCR_AUTOCOMP_DAYS, 10)

			setBoolFn(config, cmdutil.CONF_ALLOW_ARCHIVED_TAGS, true, false)

			setBoolFn(config, cmdutil.CONF_TIME_ENTRY_DEFAULTS, false, true)

			config.EXPECT().Save().Once().Return(nil)

			f.EXPECT().UI().Return(ui.NewUI(in, out, out))

			_, err := ini.NewCmdInit(f).ExecuteC()
			return err
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("Token:")
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

			c.ExpectString(
				"Look for default parameters for time entries per folder?")
			c.SendLine("?")
			c.ExpectString("https://")
			c.SendLine("y")
			c.ExpectString("Yes")

			c.ExpectEOF()
		})
}

func TestInitCmdCtrlC(t *testing.T) {
	consoletest.RunTestConsole(t,
		func(out consoletest.FileWriter, in consoletest.FileReader) error {
			f := mocks.NewMockFactory(t)
			config := mocks.NewMockConfig(t)

			f.EXPECT().Config().Return(config)
			config.EXPECT().GetString(cmdutil.CONF_TOKEN).Return("")

			f.EXPECT().UI().Return(ui.NewUI(in, out, out))

			_, err := ini.NewCmdInit(f).ExecuteC()
			if !assert.Error(t, err) {
				return errors.New("should have failed")
			}

			assert.ErrorIs(t, err, terminal.InterruptErr)
			return nil
		},
		func(c consoletest.ExpectConsole) {
			c.ExpectString("Token: ")
			c.Send(string(terminal.KeyInterrupt))

			c.ExpectEOF()
		})
}
