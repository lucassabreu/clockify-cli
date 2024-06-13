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
	"golang.org/x/sync/errgroup"
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
	Client      string
	Projects    []string
	TagIDs      []string
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
		OutputFlags: util.OutputFlags{TimeFormat: timehlp.FullTimeFormat, TimeZone: "local"},
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
	cmd.Flags().StringSliceVarP(&rf.Projects, "project", "p", []string{},
		"Will filter time entries using this project")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f, f.Config()))
	cmd.Flags().StringVarP(&rf.Client, "client", "c", "",
		"Will filter projects from this client")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f, f.Config()))
	cmd.Flags().StringSliceVarP(&rf.TagIDs, "tag", "T", []string{},
		"Will filter time entries using these tags")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "tag",
		cmdcomplutil.NewTagAutoComplete(f))

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

	cnf := f.Config()
	if len(rf.Projects) != 0 {
		if f.Config().IsAllowNameForID() {
			if rf.Projects, err = search.GetProjectsByName(
				c, cnf, workspace, rf.Client, rf.Projects); err != nil {
				return err
			}
		}
	} else if rf.Client != "" {
		if f.Config().IsAllowNameForID() {
			if rf.Client, err = search.GetClientByName(
				c, workspace, rf.Client); err != nil {
				return err
			}
		}

		ps, err := c.GetProjects(api.GetProjectsParam{
			Workspace:       workspace,
			Clients:         []string{rf.Client},
			Hydrate:         false,
			PaginationParam: api.AllPages(),
		})
		if err != nil {
			return err
		}

		rf.Projects = make([]string, len(ps))
		for i := range ps {
			rf.Projects[i] = ps[i].ID
		}
	}

	if len(rf.TagIDs) > 0 && f.Config().IsAllowNameForID() {
		if rf.TagIDs, err = search.GetTagsByName(
			c, workspace, rf.TagIDs); err != nil {
			return err
		}
	}

	if len(rf.Projects) == 0 {
		rf.Projects = []string{""}
	}

	start = timehlp.TruncateDate(start)
	end = timehlp.TruncateDate(end).Add(time.Hour * 24)

	wg := errgroup.Group{}
	logs := make([][]dto.TimeEntry, len(rf.Projects))

	for i := range rf.Projects {
		i := i
		wg.Go(func() error {
			var err error
			logs[i], err = c.LogRange(api.LogRangeParam{
				Workspace:       workspace,
				UserID:          userId,
				FirstDate:       start,
				LastDate:        end,
				Description:     rf.Description,
				ProjectID:       rf.Projects[i],
				TagIDs:          rf.TagIDs,
				PaginationParam: api.AllPages(),
			})

			return err
		})
	}

	if err = wg.Wait(); err != nil {
		return err
	}

	log := make([]dto.TimeEntry, 0)
	for i := range logs {
		log = append(log, logs[i]...)
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
		l := log
		log = make([]dto.TimeEntry, 0, len(l))
		log = append(log, fillMissing(start, l[0].TimeInterval.Start)...)

		nextDay := start
		for i := range l {
			log = append(log,
				fillMissing(nextDay, l[i].TimeInterval.Start)...)
			log = append(log, l[i])
			nextDay = l[i].TimeInterval.Start.Add(
				time.Duration(24-l[i].TimeInterval.Start.Hour()) * time.Hour)
		}

		log = append(log, fillMissing(nextDay, end)...)
	}

	return util.PrintTimeEntries(
		log, out, cnf, rf.OutputFlags)
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
