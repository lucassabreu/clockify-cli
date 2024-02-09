package mocks

import (
	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

type Factory interface {
	cmdutil.Factory
}

type Config interface {
	cmdutil.Config
}

type Client interface {
	api.Client
}

//go:generate mockery --name=TimeEntryDefaults --inpackage --with-expecter
type TimeEntryDefaults interface {
	defaults.TimeEntryDefaults
}
