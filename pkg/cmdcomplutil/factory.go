package cmdcomplutil

import "github.com/lucassabreu/clockify-cli/api"

type config interface {
	IsAllowNameForID() bool
	IsSearchProjectWithClientsName() bool
}

type factory interface {
	Client() (api.Client, error)
	GetWorkspaceID() (string, error)
}
