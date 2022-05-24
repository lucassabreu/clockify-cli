package project

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ProjectPrintWithTemplate will print each worspace using the format string
func ProjectPrintWithTemplate(format string) func([]dto.Project, io.Writer) error {
	return func(ps []dto.Project, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ps); i++ {
			if err := t.Execute(w, ps[i]); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
