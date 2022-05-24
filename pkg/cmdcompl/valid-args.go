package cmdcompl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ValidArgs interface {
	IntoUse() string
	OnlyArgs() []string
	IntoValidArgs() []string
}

// EmptyValidArgs returns a ValidArgs with no options
func EmptyValidArgs() ValidArgs {
	return new(ValidArgsSlide)
}

type ValidArgsMap map[string]string

func (va ValidArgsMap) Set(k, v string) ValidArgsMap {
	va[k] = v
	return va
}

func (va ValidArgsMap) OnlyArgs() []string {
	keys := make([]string, len(va))
	i := 0
	for k := range va {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

func (va ValidArgsMap) IntoUse() string {
	return "[" + strings.Join(va.OnlyArgs(), "|") + "]"
}

func (va ValidArgsMap) IntoValidArgs() []string {
	var args []string
	for k, v := range va {
		args = append(args, k+"\t"+v)
	}
	return args
}

func (va ValidArgsMap) Long() string {
	lenName := 0
	for k := range va {
		if len(k) > lenName {
			lenName = len(k)
		}
	}

	ft := "\t%-" + strconv.Itoa(lenName) + "s\t%s\n"
	str := ""
	for _, k := range va.OnlyArgs() {
		v := va[k]
		str = str + fmt.Sprintf(ft, k, v)
	}

	return str
}

type ValidArgsSlide []string

func (va ValidArgsSlide) IntoUse() string {
	return "[" + strings.Join(va, "|") + "]"
}

func (va ValidArgsSlide) IntoValidArgs() []string {
	return va
}

func (va ValidArgsSlide) OnlyArgs() []string {
	return va
}
