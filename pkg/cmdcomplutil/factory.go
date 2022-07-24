package cmdcomplutil

import "github.com/lucassabreu/clockify-cli/api"

type factory interface {
	Client() (api.Client, error)
	GetWorkspaceID() (string, error)
}
