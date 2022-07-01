package timeentry

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/clone"
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
	cmds = append(cmds, in.NewCmdIn(f))
	cmds = append(cmds, manual.NewCmdManual(f))
	cmds = append(cmds, clone.NewCmdClone(f))

	cmds = append(cmds, edit.NewCmdEdit(f))
	cmds = append(cmds, em.NewCmdEditMultiple(f))

	cmds = append(cmds, out.NewCmdOut(f))

	cmds = append(cmds, invoiced.NewCmdInvoiced(f)...)

	cmds = append(cmds, del.NewCmdDelete(f))

	cmds = append(cmds, show.NewCmdShow(f))
	cmds = append(cmds, report.NewCmdReport(f))

	return
}
