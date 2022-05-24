package user

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// UserPrintWithTemplate will print each worspace using the format string
func UserPrintWithTemplate(format string) func([]dto.User, io.Writer) error {
	return func(users []dto.User, w io.Writer) error {
		t, err := template.New("tmpl").Parse(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(users); i++ {
			if err := t.Execute(w, users[i]); err != nil {
				return err
			}
			fmt.Fprintln(w)
		}
		return nil
	}
}
