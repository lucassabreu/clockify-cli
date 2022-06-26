package workspace

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// WorkspacePrintWithTemplate will print each worspace using the format string
func WorkspacePrintWithTemplate(
	format string) func([]dto.Workspace, io.Writer) error {
	return func(ws []dto.Workspace, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(ws); i++ {
			if err := t.Execute(w, ws[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
