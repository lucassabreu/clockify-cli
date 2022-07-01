package cmdcompl

import (
	"sort"
	"strings"
)

type ValidArgs interface {
	// IntoUse will return a string with a complete arg use
	// Example: "{ arg1 | arg2 | arg3 }"
	IntoUse() string
	// IntoUse will return a string with the joined options
	// Example: "arg1 | arg2 | arg3"
	IntoUseOptions() string
	// OnlyArgs will return a []string to be used on cobra.Command.ValidArgs
	OnlyArgs() []string
	// OnlyArgs will return a []string to be used as result for auto-complete
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

func (va ValidArgsMap) IntoUseOptions() string {
	return strings.Join(va.OnlyArgs(), " | ")
}

func (va ValidArgsMap) IntoUse() string {
	return "{ " + va.IntoUseOptions() + " }"
}

func (va ValidArgsMap) IntoValidArgs() []string {
	var args []string
	for k, v := range va {
		args = append(args, k+"\t"+v)
	}
	return args
}

func (va ValidArgsMap) Long() string {
	str := ""
	for _, k := range va.OnlyArgs() {
		str = str + " - " + k + ": " + va[k] + "\n"
	}

	return str
}

type ValidArgsSlide []string

func (va ValidArgsSlide) IntoUseOptions() string {
	return strings.Join(va, " | ")
}

func (va ValidArgsSlide) IntoUse() string {
	return "{ " + va.IntoUseOptions() + " }"
}

func (va ValidArgsSlide) IntoValidArgs() []string {
	return va
}

func (va ValidArgsSlide) OnlyArgs() []string {
	return va
}
