package output

import (
	"io"

	"github.com/lucassabreu/clockify-cli/reportsapi"
)

func SummaryPrintTable(sr reportsapi.SummaryReport, w io.Writer) error {
	var fn func(g []reportsapi.SummaryGroup) int
	fn = func(g []reportsapi.SummaryGroup) int {
		if len(g) == 0 {
			return 0
		}

		return fn(g[0].Children) + 1
	}

	// gCount := fn(sr.GroupOne)

	return TotalsPrintTable(sr.Totals, w)
}
