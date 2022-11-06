package strhlp_test

import (
	"strings"
	"testing"

	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tts := map[string]string{
		"some long string": "Some Long STRING",
		"atencao":          "Atenção",
		"it should've this kind of keep \"stuff\"": "It should've this kind of keep \"STUFF\"",
	}

	for expected, argument := range tts {
		t.Run(expected, func(t *testing.T) {
			assert.Equal(t, expected, strhlp.Normalize(argument))
		})
	}
}
func TestInSlice(t *testing.T) {
	tts := []struct {
		name   string
		b      bool
		search string
		list   []string
	}{
		{
			name:   "unique",
			b:      true,
			search: "str 3",
			list:   []string{"str 0", "str 1", "str 2", "str 3", "str 4"},
		},
		{
			name:   "first",
			b:      true,
			search: "str 1",
			list:   []string{"str 0", "str 1", "str 1", "str 1", "str 2"},
		},
		{
			name:   "unordered",
			b:      true,
			search: "str 1",
			list:   []string{"str 0", "str 3", "str 4", "str 2", "str 1"},
		},
		{
			name:   "not found",
			b:      false,
			search: "str a",
			list:   []string{"str 0", "str 3", "str 4", "str 2", "str 1"},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.b, strhlp.InSlice(
				tt.search,
				tt.list,
			))
		})
	}
}

func TestSearch(t *testing.T) {
	tts := []struct {
		name   string
		pos    int
		search string
		list   []string
	}{
		{
			name:   "unique",
			pos:    3,
			search: "str 3",
			list:   []string{"str 0", "str 1", "str 2", "str 3", "str 4"},
		},
		{
			name:   "first",
			pos:    1,
			search: "str 1",
			list:   []string{"str 0", "str 1", "str 1", "str 1", "str 2"},
		},
		{
			name:   "unordered",
			pos:    4,
			search: "str 1",
			list:   []string{"str 0", "str 3", "str 4", "str 2", "str 1"},
		},
		{
			name:   "not found",
			pos:    -1,
			search: "str a",
			list:   []string{"str 0", "str 3", "str 4", "str 2", "str 1"},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.pos, strhlp.Search(
				tt.search,
				tt.list,
			))
		})
	}
}
func TestMap(t *testing.T) {
	tts := []struct {
		name     string
		fn       func(string) string
		expected []string
	}{
		{
			name:     "upper",
			fn:       strings.ToUpper,
			expected: []string{"VALUE 1", "VALUE 2", "VALUE 3"},
		},
		{
			name:     "last digit",
			fn:       func(s string) string { return s[len(s)-1:] },
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "same",
			fn:       func(s string) string { return "same" },
			expected: []string{"same", "same", "same"},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, strhlp.Map(
				tt.fn,
				[]string{"value 1", "value 2", "value 3"},
			))
		})
	}
}

func TestFilter(t *testing.T) {
	tts := []struct {
		name     string
		fn       func(string) bool
		expected []string
	}{
		{
			name:     "keep all",
			fn:       func(s string) bool { return true },
			expected: []string{"first", "second", "third"},
		},
		{
			name:     "keep none",
			fn:       func(s string) bool { return false },
			expected: []string{},
		},
		{
			name:     "keep second",
			fn:       func(s string) bool { return s == "second" },
			expected: []string{"second"},
		},
		{
			name:     "keep all but second",
			fn:       func(s string) bool { return s != "second" },
			expected: []string{"first", "third"},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, strhlp.Filter(
				tt.fn,
				[]string{"first", "second", "third"},
			))
		})
	}
}

func TestUnique(t *testing.T) {
	tts := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "one",
			args:     []string{"one", "one", "one"},
			expected: []string{"one"},
		},
		{
			name:     "ordered",
			args:     []string{"1", "1", "2", "2", "3", "3"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "shuffled",
			args:     []string{"2", "3", "2", "1", "3", "2"},
			expected: []string{"2", "3", "1"},
		},
		{
			name:     "no changes",
			args:     []string{"2", "3", "4", "1", "5", "6"},
			expected: []string{"2", "3", "4", "1", "5", "6"},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, strhlp.Unique(tt.args))
		})
	}
}

func TestPadSpace(t *testing.T) {
	tts := []struct {
		name     string
		size     int
		expected string
	}{
		{
			name:     "zero",
			size:     0,
			expected: "some word",
		},
		{
			name:     "same size",
			size:     9,
			expected: "some word",
		},
		{
			name:     "10",
			size:     10,
			expected: "some word ",
		},
		{
			name:     "30",
			size:     30,
			expected: "some word                     ",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, strhlp.PadSpace(
				"some word",
				tt.size,
			))
		})
	}
}

func TestListForHumans(t *testing.T) {
	tts := []struct {
		args     []string
		expected string
	}{
		{
			args:     []string{"uno", "dos", "tres"},
			expected: "uno, dos and tres",
		},
		{
			args:     []string{"first", "last"},
			expected: "first and last",
		},
		{
			args:     []string{"1", "2", "3", "go!"},
			expected: "1, 2, 3 and go!",
		},
		{
			args:     []string{"only one"},
			expected: "only one",
		},
	}

	for _, tt := range tts {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, strhlp.ListForHumans(tt.args))
		})
	}
}
