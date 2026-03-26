package user

import (
	"fmt"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// UserPrintQuietly will only print the IDs
func UserPrintQuietly(users []dto.User, w io.Writer) error {
	for i := 0; i < len(users); i++ {
		if _, err := fmt.Fprintln(w, users[i].ID); err != nil {
			return err
		}
	}

	return nil
}
