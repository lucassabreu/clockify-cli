package output

import (
	"io"

	"github.com/lucassabreu/clockify-cli/reportsapi"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var messagePrinter = message.NewPrinter(language.English)

func inttostr(i int) string {
	return messagePrinter.Sprint("%d", i)
}
func dectostr(i float32) string {
	return messagePrinter.Sprint("%.2f", i)
}

func TotalsPrintTable(ts []reportsapi.Total, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.SetHeader([]string{"Total Time", "Total Billable Time", "Entries Count", "Total Amount"})

	lines := make([][]string, len(ts))
	for i, t := range ts {
		lines[i] = []string{
			inttostr(t.TotalTime),
			inttostr(t.TotalBillableTime),
			inttostr(t.EntriesCount),
			dectostr(t.TotalAmount),
		}
	}

	tw.AppendBulk(lines)
	tw.Render()

	return nil
}
