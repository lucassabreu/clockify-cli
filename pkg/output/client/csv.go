package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ClientsCSVPrint will print as CSV
func ClientsCSVPrint(clients []dto.Client, out io.Writer) error {
	w := csv.NewWriter(out)

	if err := w.Write([]string{
		"id",
		"name",
		"archived",
	}); err != nil {
		return err
	}

	for i := 0; i < len(clients); i++ {
		c := clients[i]
		if err := w.Write([]string{
			c.ID,
			c.Name,
			fmt.Sprintf("%v", c.Archived),
		}); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}
