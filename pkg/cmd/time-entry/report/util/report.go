package util

import (
	"io"
	"sort"
	"time"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/lucassabreu/clockify-cli/pkg/cmdcomplutil"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
	"github.com/lucassabreu/clockify-cli/pkg/search"
	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/spf13/cobra"
)

const (
	HelpNamesForIds           = util.HelpNamesForIds
	HelpMoreInfoAboutPrinting = util.HelpMoreInfoAboutPrinting
)

// ReportFlags reads the "shared" flags for report commands
type ReportFlags struct {
	util.OutputFlags

	FillMissingDates bool

	Billable    bool
	NotBillable bool

	Description string
	Project     string
}

// Check will assure that there is no conflicting flag values
func (rf ReportFlags) Check() error {
	if err := rf.OutputFlags.Check(); err != nil {
		return err
	}

	return cmdutil.XorFlag(map[string]bool{
		"billable":     rf.Billable,
		"not-billable": rf.NotBillable,
	})
}

// NewReportFlags helps creating a util.ReportFlags for report commands
func NewReportFlags() ReportFlags {
	return ReportFlags{
		OutputFlags: util.OutputFlags{TimeFormat: timehlp.FullTimeFormat},
	}
}

// AddReportFlags add flags for print out the time entries
func AddReportFlags(
	f cmdutil.Factory, cmd *cobra.Command, rf *ReportFlags,
) {
	util.AddPrintTimeEntriesFlags(cmd, &rf.OutputFlags)
	util.AddPrintMultipleTimeEntriesFlags(cmd)

	cmd.Flags().BoolVarP(&rf.FillMissingDates, "fill-missing-dates", "e", false,
		"add empty lines for dates without time entries")
	cmd.Flags().StringVarP(&rf.Description, "description", "d", "",
		"will filter time entries that contains this on the description field")
	cmd.Flags().StringVarP(&rf.Project, "project", "p", "",
		"Will filter time entries using this project")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f))

	cmd.Flags().BoolVar(&rf.Billable, "billable", false,
		"Will filter time entries that are billable")
	cmd.Flags().BoolVar(&rf.NotBillable, "not-billable", false,
		"Will filter time entries that are not billable")
}

// ReportWithRange fetches and prints out time entries
func ReportWithRange(
	f cmdutil.Factory, start, end time.Time,
	out io.Writer, rf ReportFlags,
) error {
	userId, err := f.GetUserID()
	if err != nil {
		return err
	}

	workspace, err := f.GetWorkspaceID()
	if err != nil {
		return err
	}

	c, err := f.Client()
	if err != nil {
		return err
	}

	if rf.Project != "" && f.Config().IsAllowNameForID() {
		if rf.Project, err = search.GetProjectByName(
			c, workspace, rf.Project); err != nil {
			return err
		}
	}

	start = timehlp.TruncateDate(start)
	end = timehlp.TruncateDate(end).Add(time.Hour * 24)
	log, err := c.LogRange(api.LogRangeParam{
		Workspace:       workspace,
		UserID:          userId,
		FirstDate:       start,
		LastDate:        end,
		Description:     rf.Description,
		ProjectID:       rf.Project,
		PaginationParam: api.AllPages(),
	})

	if err != nil {
		return err
	}

	if rf.Billable || rf.NotBillable {
		log = filterBilling(log, rf.Billable)
	}

	sort.Slice(log, func(i, j int) bool {
		return log[j].TimeInterval.Start.After(
			log[i].TimeInterval.Start,
		)
	})

	if rf.FillMissingDates && len(log) > 0 {
		newLog := make([]dto.TimeEntry, 0, len(log))

		newLog = append(newLog,
			fillMissing(start, log[0].TimeInterval.Start)...)
		nextDay := start
		for _, t := range log {
			newLog = append(newLog,
				fillMissing(nextDay, t.TimeInterval.Start)...)
			newLog = append(newLog, t)
			nextDay = t.TimeInterval.Start.Add(
				time.Duration(24-t.TimeInterval.Start.Hour()) * time.Hour)
		}
		log = append(newLog, fillMissing(nextDay, end)...)
	}

	return util.PrintTimeEntries(
		log, out, f.Config(), rf.OutputFlags)
}

func filterBilling(l []dto.TimeEntry, billable bool) []dto.TimeEntry {
	r := make([]dto.TimeEntry, 0, len(l))
	for i := 0; i < len(l); i++ {
		if l[i].Billable == billable {
			r = append(r, l[i])
		}
	}

	return r
}

func fillMissing(first, last time.Time) []dto.TimeEntry {
	first = timehlp.TruncateDate(first)
	last = timehlp.TruncateDate(last)

	d := int(last.Sub(first).Hours() / 24)
	if d <= 0 {
		return []dto.TimeEntry{}
	}

	missing := make([]dto.TimeEntry, d)
	for i := 0; i < d; i++ {
		t := missing[i]
		ti := first.AddDate(0, 0, i)
		t.TimeInterval.Start = ti
		t.TimeInterval.End = &ti
		missing[i] = t
	}

	return missing
}
