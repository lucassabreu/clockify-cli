package task

import (
	"encoding/json"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TasksJSONPrint will print as JSON
func TasksJSONPrint(t []dto.Task, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}
