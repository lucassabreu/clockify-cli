package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/lucassabreu/clockify-cli/api/dto"
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
			return timehlp.Now().UTC()
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
	"ident": func(s, prefix string) string {
		return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix)
	},
	"since": func(s time.Time, e ...time.Time) dto.Duration {
		return diff(s, firstOrNow(e))
	},
	"until": func(s time.Time, e ...time.Time) dto.Duration {
		return diff(firstOrNow(e), s)
	},
}

func firstOrNow(ts []time.Time) time.Time {
	if len(ts) == 0 {
		return timehlp.Now().UTC()
	}
	return ts[0]
}

func diff(s, e time.Time) dto.Duration {
	return dto.Duration{Duration: e.Sub(s)}
}

func NewTemplate(format string) (*template.Template, error) {
	format = strings.ReplaceAll(format, "\\n", "\n")
	format = strings.ReplaceAll(format, "\\t", "\t")
	return template.New("tmpl").Funcs(funcMap).Parse(format + "\n")
}
