package search

import (
	"fmt"
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

func TestFindByNameAmbiguous(t *testing.T) {
	entities := []named{
		namedStruct{ID: "1", Name: "Web Redesign"},
		namedStruct{ID: "2", Name: "Web Maintenance"},
		namedStruct{ID: "3", Name: "Backend"},
	}

	id, err := findByName("web", "project", func() ([]named, error) {
		return entities, nil
	})

	assert.Error(t, err)
	assert.Equal(t, "web", id)

	var eAmbiguous ErrAmbiguous
	if assert.ErrorAs(t, err, &eAmbiguous) {
		assert.Equal(t, "project", eAmbiguous.EntityName)
		assert.Equal(t, "web", eAmbiguous.Reference)
		if assert.Len(t, eAmbiguous.Candidates, 2) {
			assert.Equal(t, Candidate{ID: "2", Name: "Web Maintenance"},
				eAmbiguous.Candidates[0])
			assert.Equal(t, Candidate{ID: "1", Name: "Web Redesign"},
				eAmbiguous.Candidates[1])
		}
	}

	assert.Contains(t, err.Error(), "Web Redesign")
	assert.Contains(t, err.Error(), "Web Maintenance")
}

func TestFindByNameAmbiguousCapsCandidatesInMessage(t *testing.T) {
	entities := make([]named, 0, 12)
	for i := 1; i <= 12; i++ {
		entities = append(entities, namedStruct{
			ID:   fmt.Sprintf("p%d", i),
			Name: fmt.Sprintf("Web %02d", i),
		})
	}

	_, err := findByName("web", "project", func() ([]named, error) {
		return entities, nil
	})

	if assert.Error(t, err) {
		msg := err.Error()
		assert.Contains(t, msg, "'Web 10' (p10)")
		assert.NotContains(t, msg, "'Web 11'")
		assert.Contains(t, msg, " and 2 more")
		assert.NotContains(t, msg, "and and")
	}
}

func TestFindByNameAmbiguousSortsCandidates(t *testing.T) {
	redesign := namedStruct{ID: "1", Name: "Web Redesign"}
	maintenance := namedStruct{ID: "2", Name: "Web Maintenance"}
	design := namedStruct{ID: "3", Name: "Web Design"}

	expected := []Candidate{
		{ID: "3", Name: "Web Design"},
		{ID: "2", Name: "Web Maintenance"},
		{ID: "1", Name: "Web Redesign"},
	}

	for _, entities := range [][]named{
		{redesign, maintenance, design},
		{design, redesign, maintenance},
		{maintenance, design, redesign},
	} {
		_, err := findByName("web", "project", func() ([]named, error) {
			return entities, nil
		})

		var eAmbiguous ErrAmbiguous
		if assert.ErrorAs(t, err, &eAmbiguous) {
			assert.Equal(t, expected, eAmbiguous.Candidates)
		}
	}
}

func TestFindByNameUnambiguous(t *testing.T) {
	entities := []named{
		namedStruct{ID: "1", Name: "Web Redesign"},
		namedStruct{ID: "2", Name: "Web Maintenance"},
		namedStruct{ID: "3", Name: "Backend"},
		namedStruct{ID: "web-4", Name: "Something else"},
		namedStruct{ID: "w", Name: "web"},
	}

	find := func(r string) (string, error) {
		return findByName(r, "project", func() ([]named, error) {
			return entities, nil
		})
	}

	t.Run("exact id", func(t *testing.T) {
		id, err := find("web-4")
		if assert.NoError(t, err) {
			assert.Equal(t, "web-4", id)
		}
	})

	t.Run("exact name wins over other fuzzy matches", func(t *testing.T) {
		id, err := find("web")
		if assert.NoError(t, err) {
			assert.Equal(t, "w", id)
		}
	})

	t.Run("single fuzzy match", func(t *testing.T) {
		id, err := find("maintenance")
		if assert.NoError(t, err) {
			assert.Equal(t, "2", id)
		}
	})

	t.Run("repeated ids are not ambiguous", func(t *testing.T) {
		repeated := []named{
			namedStruct{ID: "1", Name: "Web Redesign"},
			namedStruct{ID: "c", Name: "clockify"},
			namedStruct{ID: "2", Name: "Web Maintenance"},
			namedStruct{ID: "c", Name: "clockify"},
		}

		id, err := findByName("clockify", "client", func() ([]named, error) {
			return repeated, nil
		})
		if assert.NoError(t, err) {
			assert.Equal(t, "c", id)
		}
	})
}
