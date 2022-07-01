package tag

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// TagPrintWithTemplate will print each worspace using the format string
func TagPrintWithTemplate(format string) func([]dto.Tag, io.Writer) error {
	return func(ts []dto.Tag, w io.Writer) error {
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
