package timeentry

import (
	_ "embed"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

//go:embed template.gotmpl.md
var mdTemplate string

// TimeEntriesMarkdownPrint will print time entries in "markdown blocks"
func TimeEntriesMarkdownPrint(tes []dto.TimeEntry, w io.Writer) error {
	return TimeEntriesPrintWithTemplate(mdTemplate)(tes, w)
}
