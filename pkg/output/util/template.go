package util

import (
	"bytes"
	"encoding/json"
	"text/template"
	"time"

	"github.com/lucassabreu/clockify-cli/pkg/timehlp"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"gopkg.in/yaml.v3"
)

func formatTime(f string) func(time.Time) string {
	return func(t time.Time) string {
		return t.Format(f)
	}
}

var funcMap = template.FuncMap{
	"formatDateTime": formatTime(timehlp.FullTimeFormat),
	"fdt":            formatTime(timehlp.FullTimeFormat),
	"formatTime":     formatTime(timehlp.OnlyTimeFormat),
	"ft":             formatTime(timehlp.OnlyTimeFormat),
	"now": func(t *time.Time) time.Time {
		if t == nil {
			return time.Now().UTC()
		}

		return *t
	},
	"json": func(j interface{}) string {
		w := bytes.NewBufferString("")
		if err := json.NewEncoder(w).Encode(j); err != nil {
			return ""
		}

		return w.String()
	},
	"yaml": func(j interface{}) string {
		w := bytes.NewBufferString("")
		if err := yaml.NewEncoder(w).Encode(j); err != nil {
			return ""
		}

		return w.String()
	},
	"pad": strhlp.PadSpace,
}

func NewTemplate(format string) (*template.Template, error) {
	return template.New("tmpl").Funcs(funcMap).Parse(format + "\n")
}
