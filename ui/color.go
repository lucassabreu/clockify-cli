package ui

import (
	"errors"
	"strconv"
)

func HEX(hex string) (RGB, error) {
	if len(hex) != 6 {
		return RGB{}, errors.New("HEX must be 6 characters")
	}

	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	return RGB{
		int(values >> 16),
		int((values >> 8) & 0xFF),
		int(values & 0xFF),
	}, nil
}

type RGB [3]int

func (c RGB) R() int {
	return c[0]
}

func (c RGB) G() int {
	return c[1]
}

func (c RGB) B() int {
	return c[2]
}

func (c RGB) Values() []int {
	return []int{c.R(), c.G(), c.B()}
}
