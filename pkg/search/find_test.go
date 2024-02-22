package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchOnList(t *testing.T) {
	entities := []named{
		namedStruct{
			ID:   "1",
			Name: "entity one",
		},
		namedStruct{
			ID:   "2",
			Name: "entity two",
		},
		namedStruct{
			ID:   "3",
			Name: "entity three",
		},
		namedStruct{
			ID:   "4",
			Name: "more complex name",
		},
		namedStruct{
			ID:   "id",
			Name: "by id",
		},
		namedStruct{
			ID:   "bra",
			Name: "with [bracket]",
		},
	}

	tts := []struct {
		name     string
		search   string
		entities []named
		result   string
	}{
		{
			name:     "one term",
			search:   "two",
			entities: entities,
			result:   "2",
		},
		{
			name:     "two terms",
			search:   "complex name",
			entities: entities,
			result:   "4",
		},
		{
			name:     "sections of the name",
			search:   "mo nam",
			entities: entities,
			result:   "4",
		},
		{
			name:     "with brackets",
			search:   "[bracket]",
			entities: entities,
			result:   "bra",
		},
		{
			name:     "using id",
			search:   "by id",
			entities: entities,
			result:   "id",
		},
	}

	for i := range tts {
		tt := tts[i]
		t.Run(tt.name, func(t *testing.T) {
			id, err := findByName(tt.search, "element", func() ([]named, error) {
				return tt.entities, nil
			})

			if assert.NoError(t, err) {
				assert.Equal(t, tt.result, id)
			}
		})
	}
}
