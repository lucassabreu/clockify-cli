package completion

import (
	"strings"
)

type ValidArgs interface {
	IntoUse() string
	OnlyArgs() []string
	IntoValidArgs() []string
}

// EmptyValidArgs returns a ValidArgs with no options
func EmptyValidArgs() ValidArgs {
	return new(ValigsArgsSlide)
}

type ValigsArgsMap map[string]string

func (va ValigsArgsMap) Set(k, v string) ValigsArgsMap {
	va[k] = v
	return va
}

func (va ValigsArgsMap) OnlyArgs() []string {
	var keys []string
	for k := range va {
		keys = append(keys, k)
	}
	return keys
}

func (va ValigsArgsMap) IntoUse() string {

	return "[" + strings.Join(va.OnlyArgs(), "|") + "]"
}

func (va ValigsArgsMap) IntoValidArgs() []string {
	var args []string
	for k, v := range va {
		args = append(args, k+"\t"+v)
	}
	return args
}

type ValigsArgsSlide []string

func (va ValigsArgsSlide) IntoUse() string {
	return "[" + strings.Join(va, "|") + "]"
}

func (va ValigsArgsSlide) IntoValidArgs() []string {
	return va
}

func (va ValigsArgsSlide) OnlyArgs() []string {
	return va
}
