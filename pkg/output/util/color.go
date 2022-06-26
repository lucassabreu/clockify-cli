package util

import (
	"os"

	"github.com/lucassabreu/clockify-cli/pkg/ui"
)

// ColorToTermColor coverts HEX color to term colors
func ColorToTermColor(hex string) []int {
	if hex == "" {
		return []int{}
	}

	fi, _ := os.Stdout.Stat()
	if fi.Mode()&os.ModeCharDevice == 0 {
		return []int{}
	}

	if c, err := ui.HEX(hex[1:]); err == nil {
		return append(
			[]int{38, 2},
			c.Values()...,
		)
	}

	return []int{}
}
