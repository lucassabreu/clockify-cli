package project

import (
	"encoding/json"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ProjectsJSONPrint will print as JSON
func ProjectsJSONPrint(t []dto.Project, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// ProjectJSONPrint will print as JSON
func ProjectJSONPrint(t dto.Project, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}
