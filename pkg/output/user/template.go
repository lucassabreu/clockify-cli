package user

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// UserPrintWithTemplate will print each worspace using the format string
func UserPrintWithTemplate(format string) func([]dto.User, io.Writer) error {
	return func(users []dto.User, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(users); i++ {
			if err := t.Execute(w, users[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
