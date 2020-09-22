package output

import (
	"io"

	"github.com/lucassabreu/clockify-cli/reportsapi"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/olekukonko/tablewriter"
)

func SummaryPrintTable(gs reportsapi.GroupSlice, sr reportsapi.SummaryReport, w io.Writer) error {
	tw := tablewriter.NewWriter(w)

	headers := make([]string, 0)
	for _, g := range gs {
		headers = append(headers, string(g))
	}
	headers = append(headers, "Duration", "Amount")

	tw.SetHeader(headers)
	tw.SetRowLine(true)

	tw.AppendBulk(recurSummaryTable(sr.GroupOne, len(gs), 0))

	tw.Render()
	return TotalsPrintTable(sr.Totals, w)
}

func recurSummaryTable(gs []reportsapi.SummaryGroup, groupCount, current int) (result [][]string) {
	labels := make([]string, groupCount)
	for _, g := range gs {
		labels[current] = g.Name
		result = append(result, strhlp.Merge(
			labels,
			[]string{
				inttostr(g.Duration),
				dectostr(g.Amount),
			},
		))

		if len(g.Children) > 0 {
			result = append(result, recurSummaryTable(g.Children, groupCount, current+1)...)
		}
	}

	return result
}
