package util_test

import (
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/report/util"
	"github.com/stretchr/testify/assert"
)

func TestReportFlags_Check(t *testing.T) {
	tts := map[string]struct {
		rf  util.ReportFlags
		err string
	}{
		"just billable": {
			rf: util.ReportFlags{
				Billable:    true,
				NotBillable: false,
			},
		},
		"just not-billable": {
			rf: util.ReportFlags{
				Billable:    false,
				NotBillable: true,
			},
		},
		"only one billable": {
			rf: util.ReportFlags{
				Billable:    true,
				NotBillable: true,
			},
			err: "can't be used together.*billable.*not-billable",
		},
		"just client": {
			rf: util.ReportFlags{
				Client: "me",
			},
		},
		"just projects": {
			rf: util.ReportFlags{
				Projects: []string{"mine"},
			},
		},
		"client and project": {
			rf: util.ReportFlags{
				Client:   "me",
				Projects: []string{"mine"},
			},
		},
		"fill missing dates": {
			rf: util.ReportFlags{
				FillMissingDates: true,
			},
		},
		"limit": {
			rf: util.ReportFlags{
				Limit: 10,
			},
		},
		"only limit or fill missing": {
			rf: util.ReportFlags{
				Limit:            10,
				FillMissingDates: true,
			},
			err: "can't be used together.*fill-missing-dates.*limit",
		},
		"limit and page": {
			rf: util.ReportFlags{
				Limit: 10,
				Page:  10,
			},
		},
		"page needs limit": {
			rf: util.ReportFlags{
				Page: 10,
			},
			err: "page can't be used without limit",
		},
	}

	for name, tt := range tts {
		t.Run(name, func(t *testing.T) {
			err := tt.rf.Check()

			if tt.err == "" {
				assert.NoError(t, err)
				return
			}

			if !assert.Error(t, err) {
				return
			}

			assert.Regexp(t, tt.err, err.Error())
		})
	}
}
