package timeentry

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
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
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/split"
	teutil "github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTimeEntry(f cmdutil.Factory) (cmds []*cobra.Command) {
	rmFn := func(
		tes []dto.TimeEntry, out io.Writer, of teutil.OutputFlags) error {
		return teutil.PrintTimeEntries(tes, out, f.Config(), of)
	}

	rFn := func(
		tei dto.TimeEntryImpl, out io.Writer, of teutil.OutputFlags,
	) error {
		return teutil.PrintTimeEntryImpl(tei, f, out, of)
	}

	cmds = append(
		cmds,

		in.NewCmdIn(f, rFn),
		manual.NewCmdManual(f),
		clone.NewCmdClone(f),

		edit.NewCmdEdit(f, rFn),
		em.NewCmdEditMultiple(f),

		split.NewCmdSplit(f, rmFn),

		out.NewCmdOut(f),

		del.NewCmdDelete(f),

		show.NewCmdShow(f),
		report.NewCmdReport(f),
	)

	cmds = append(cmds, invoiced.NewCmdInvoiced(f)...)

	return
}
