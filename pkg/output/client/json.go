package client

import (
	"encoding/json"
	"io"

	"github.com/lucassabreu/clockify-cli/api/dto"
)

// ClientJSONPrint will print as JSON
func ClientJSONPrint(t dto.Client, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

// ClientsJSONPrint will print as JSON
func ClientsJSONPrint(t []dto.Client, w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}
