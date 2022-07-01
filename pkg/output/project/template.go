package project

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// ProjectPrintWithTemplate will print each worspace using the format string
func ProjectPrintWithTemplate(format string) func([]dto.Project, io.Writer) error {
	return func(ps []dto.Project, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ps); i++ {
			if err := t.Execute(w, ps[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
