package cmdcompl_test

import (
	"testing"

	"github.com/lucassabreu/clockify-cli/pkg/cmdcompl"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCombineSuggestionsToArgs(t *testing.T) {
	called := make([]int, 2)
	fns := []cmdcompl.SuggestFn{
		func(_ *cobra.Command, _ []string, _ string) (cmdcompl.ValidArgs, error) {
			called[0]++
			return cmdcompl.ValidArgsSlide{"first"}, nil
		},
		func(_ *cobra.Command, _ []string, _ string) (cmdcompl.ValidArgs, error) {
			called[1]++
			return cmdcompl.ValidArgsSlide{"second"}, nil
		},
	}

	fn := cmdcompl.CombineSuggestionsToArgs(fns...)

	tts := []struct {
		name        string
		argsCount   int
		calledFn    int
		suggestions []string
	}{
		{
			name:        "no args",
			argsCount:   0,
			calledFn:    0,
			suggestions: []string{"first"},
		},
		{
			name:        "last fn",
			argsCount:   len(fns) - 1,
			calledFn:    1,
			suggestions: []string{"second"},
		},
		{
			name:        "args exhausted the fns",
			argsCount:   len(fns),
			calledFn:    -1,
			suggestions: []string{},
		},
		{
			name:        "more args than fns",
			argsCount:   len(fns) + 1,
			calledFn:    -1,
			suggestions: []string{},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			called[0], called[1] = 0, 0
			args := make([]string, tt.argsCount)
			var suggestions []string
			var directive cobra.ShellCompDirective

			assert.NotPanics(t, func() {
				suggestions, directive = fn(&cobra.Command{}, args, "")
			})
			assert.Equal(t, tt.suggestions, suggestions)
			assert.Equal(t, cobra.ShellCompDirectiveDefault, directive)

			for i, c := range called {
				if i == tt.calledFn {
					assert.Equal(t, 1, c, "fns[%d] should have been called", i)
					continue
				}
				assert.Equal(t, 0, c, "fns[%d] should not have been called", i)
			}
		})
	}
}
