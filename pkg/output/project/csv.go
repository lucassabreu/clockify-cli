package project

import (
	"encoding/csv"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ProjectsCSVPrint will print each time entry using the format string
func ProjectsCSVPrint(ps []dto.Project, out io.Writer) error {
	w := csv.NewWriter(out)

	if err := w.Write([]string{
		"id",
		"name",
		"client.id",
		"client.name",
	}); err != nil {
		return err
	}

	for i := 0; i < len(ps); i++ {
		p := ps[i]
		if err := w.Write([]string{
			p.ID,
			p.Name,
			p.ClientID,
			p.ClientName,
		}); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
