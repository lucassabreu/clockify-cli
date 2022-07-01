package user

import (
	"encoding/json"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// UserJSONPrint will print the user as a JSON
func UserJSONPrint(u dto.User, w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}
