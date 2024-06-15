package util

import (
	"bytes"
	"encoding/json"
	"slices"
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
	"formatTimeWS":   formatTime(timehlp.SimplerOnlyTimeFormat),
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
	"repeatString": strings.Repeat,
	"maxLength": func(s ...string) int {
		length := 0
		for i := range s {
			l := len(s[i])
			if l > length {
				length = l
			}
		}

		return length
	},
	"maxInt": func(s ...int) int { return slices.Max(s) },
	"concat": func(ss ...string) string {
		b := &strings.Builder{}
		for _, s := range ss {
			b.WriteString(s)
		}

		return b.String()
	},
	"dsf": func(ds string) string {
		d, err := dto.StringToDuration(ds)
		if err != nil {
			panic(err)
		}

		return dto.Duration{Duration: d}.HumanString()
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
