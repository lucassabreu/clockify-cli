package mocks

import (
	"github.com/lucassabreu/clockify-cli/pkg/cmdutil"
)

//go:generate mockery --name=Factory --inpackage
type Factory interface {
	cmdutil.Factory
}
