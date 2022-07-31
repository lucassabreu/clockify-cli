package util_test

import (
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/stretchr/testify/assert"
)

func TestReportFlagsChecks(t *testing.T) {
	rf := util.NewReportFlags()
	rf.Billable = true
	rf.NotBillable = true

	err := rf.Check()
	assert.Error(t, err)
	assert.Regexp(t,
		"can't be used together.*billable.*not-billable", err.Error())

	rf.Billable = false
	rf.NotBillable = true

	assert.NoError(t, rf.Check())

	rf.Billable = true
	rf.NotBillable = false

	assert.NoError(t, rf.Check())
}
