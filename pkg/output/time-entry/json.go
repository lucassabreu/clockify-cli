package timeentry

import (
	"encoding/json"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// TimeEntryJSONPrint will print as JSON
func TimeEntryJSONPrint(t dto.TimeEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// TimeEntriesJSONPrint will print as JSON
func TimeEntriesJSONPrint(t []dto.TimeEntry, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}
