package workspace

import (
	"fmt"
	"html/template"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// WorkspacePrintWithTemplate will print each worspace using the format string
func WorkspacePrintWithTemplate(
	format string) func([]dto.Workspace, io.Writer) error {
	return func(ws []dto.Workspace, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ws); i++ {
			if err := t.Execute(w, ws[i]); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
