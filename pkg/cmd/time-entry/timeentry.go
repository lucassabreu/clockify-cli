package timeentry

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/clone"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/defaults"
	del "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/delete"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/edit"
	em "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/edit-multipple"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/in"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/invoiced"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/manual"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/out"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/show"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTimeEntry(f cmdutil.Factory) (cmds []*cobra.Command) {
	cmds = append(
		cmds,

		in.NewCmdIn(f, nil),
		manual.NewCmdManual(f),
		clone.NewCmdClone(f),

		edit.NewCmdEdit(f, nil),
		em.NewCmdEditMultiple(f),

		out.NewCmdOut(f),

		del.NewCmdDelete(f),

		show.NewCmdShow(f),
		report.NewCmdReport(f),

		defaults.NewCmdDefaults(f),
	)

	cmds = append(cmds, invoiced.NewCmdInvoiced(f)...)

	return
}
