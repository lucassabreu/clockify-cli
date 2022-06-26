package client

import (
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/output/util"
)

// ClientPrintWithTemplate will print each client using the format string
func ClientPrintWithTemplate(format string) func([]dto.Client, io.Writer) error {
	return func(cs []dto.Client, w io.Writer) error {
		t, err := util.NewTemplate(format)
		if err != nil {
			return err
		}

		for i := 0; i < len(cs); i++ {
			if err := t.Execute(w, cs[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
