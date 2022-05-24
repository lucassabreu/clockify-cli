package tag

import (
	"fmt"
	"html/template"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TagPrintWithTemplate will print each worspace using the format string
func TagPrintWithTemplate(format string) func([]dto.Tag, io.Writer) error {
	return func(ts []dto.Tag, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ts); i++ {
			if err := t.Execute(w, ts[i]); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
