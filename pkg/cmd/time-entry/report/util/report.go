package util

import (
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

// AddReportFlags add flags for print out the time entries
func AddReportFlags(f cmdutil.Factory, cmd *cobra.Command) {
	util.AddPrintTimeEntriesFlags(cmd)
	util.AddPrintMultipleTimeEntriesFlags(cmd)

	cmd.Flags().BoolP("fill-missing-dates", "e", false,
		"add empty lines for dates without time entries")
	cmd.Flags().StringP("description", "d", "",
		"will filter time entries that contains this on the description field")
	cmd.Flags().StringP("project", "p", "",
		"Will filter time entries using this project")
	_ = cmdcompl.AddSuggestionsToFlag(cmd, "project",
		cmdcomplutil.NewProjectAutoComplete(f))
}

// ReportWithRange fetches and prints out time entries
func ReportWithRange(
	f cmdutil.Factory, start, end time.Time, cmd *cobra.Command) error {

	fillMissingDates, _ := cmd.Flags().GetBool("fill-missing-dates")
	description, _ := cmd.Flags().GetString("description")
	project, _ := cmd.Flags().GetString("project")

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

	if f.Config().IsAllowNameForID() && project != "" {
		if project, err = search.GetProjectByName(
			c, workspace, project); err != nil {
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
		Description:     description,
		ProjectID:       project,
		PaginationParam: api.AllPages(),
	})

	if err != nil {
		return err
	}

	sort.Slice(log, func(i, j int) bool {
		return log[j].TimeInterval.Start.After(
			log[i].TimeInterval.Start,
		)
	})

	if fillMissingDates && len(log) > 0 {
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

	return util.PrintTimeEntries(log, cmd, timehlp.FullTimeFormat, f.Config())
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
