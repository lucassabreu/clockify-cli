package task

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TaskPrintWithTemplate will print each client using the format string
func TaskPrintWithTemplate(format string) func([]dto.Task, io.Writer) error {
	return func(ts []dto.Task, w io.Writer) error {
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
