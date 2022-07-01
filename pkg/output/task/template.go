package task

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// TaskPrintWithTemplate will print each client using the format string
func TaskPrintWithTemplate(format string) func([]dto.Task, io.Writer) error {
	return func(ts []dto.Task, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ts); i++ {
			if err := t.Execute(w, ts[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
