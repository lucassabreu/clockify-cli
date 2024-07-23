package show_test

import (
	"bytes"
	"testing"

	outd "github.com/lucassabreu/clockify-cli/pkg/output/defaults"

	"github.com/MakeNowJust/heredoc"
	"github.com/lucassabreu/clockify-cli/internal/mocks"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults/show"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/stretchr/testify/assert"
)

var bFalse = false
var bTrue = true

func TestNewCmdShow_ShouldPrintDefaults(t *testing.T) {
	ft := func(name string,
		dte defaults.DefaultTimeEntry,
		args []string, expected string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			f := mocks.NewMockFactory(t)

			ted := mocks.NewMockTimeEntryDefaults(t)
			f.EXPECT().TimeEntryDefaults().Return(ted)
			ted.EXPECT().Read().Return(dte, nil)

			cmd := show.NewCmdShow(f, outd.Report)
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			cmd.SetArgs(args)

			out := bytes.NewBufferString("")

			cmd.SetOut(out)
			cmd.SetErr(out)

			_, err := cmd.ExecuteC()

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, expected, out.String())
		})
	}

	dte := defaults.DefaultTimeEntry{
		ProjectID: "p",
		Billable:  &bFalse,
		TagIDs:    []string{"t1"},
	}

	ft("as json", dte, []string{"--format=json"},
		`{"project":"p","billable":false,"tags":["t1"]}`)

	ft("as yaml", dte, []string{"--format=yaml"}, heredoc.Doc(`
		project: p
		billable: false
		tags: [t1]
	`))

	dte = defaults.DefaultTimeEntry{
		ProjectID: "p",
		TaskID:    "t",
		Billable:  &bTrue,
	}

	ft("as json", dte, []string{"--format=json"},
		`{"project":"p","task":"t","billable":true}`)

	ft("as yaml", dte, []string{"--format=yaml"}, heredoc.Doc(`
		project: p
		task: t
		billable: true
	`))
}
