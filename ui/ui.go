package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lucassabreu/clockify-cli/strhlp"
)

func selectFilter(filter, value string, _ int) bool {
	r := strings.Join([]string{"]", "^", `\\`, "[", ".", "(", ")", "-"}, "")
	filter = regexp.MustCompile("["+r+"]+").ReplaceAllString(strhlp.Normalize(filter), "")
	filter = regexp.MustCompile(`\s+`).ReplaceAllString(filter, " ")
	filter = strings.ReplaceAll(filter, " ", ".*")
	filter = strings.ReplaceAll(filter, "*", ".*")

	return regexp.MustCompile(filter).Match([]byte(strhlp.Normalize(value)))
}

func askString(p survey.Prompt) (string, error) {
	answer := ""
	return answer, survey.AskOne(p, &answer, nil)
}

// WithSuggestion applies the suggestion function to the input question
func WithSuggestion(fn func(toComplete string) []string) func(*survey.Input) {
	return func(i *survey.Input) {
		i.Suggest = fn
	}
}

// InputOption represets a funcion the customizes a survey.Input object
type InputOption func(*survey.Input)

// AskForText interactively ask for one string from the user
func AskForText(message, d string, opts ...InputOption) (string, error) {
	i := &survey.Input{
		Message: message,
		Default: d,
	}

	for _, o := range opts {
		o(i)
	}

	return askString(i)
}

// AskForInt interactively ask for one int from the user
func AskForInt(message string, d int) (int, error) {
	return d, survey.AskOne(
		&survey.Input{
			Message: message,
			Default: strconv.Itoa(d),
		}, &d,
		survey.WithValidator(func(ans interface{}) error {
			v, ok := ans.(string)
			if !ok {
				return fmt.Errorf("needs to be a string")
			}

			_, err := strconv.Atoi(v)
			return err
		}),
	)
}

// AskFromOptions interactively ask the user to choose one option or none
func AskFromOptions(message string, options []string, d string) (string, error) {
	p := &survey.Select{
		Message: message,
		Options: options,
		Filter:  selectFilter,
	}

	if d != "" && strhlp.Search(d, options) != -1 {
		p.Default = d
	}

	return askString(p)
}

// AskManyFromOptions interactively ask the user to choose none or many option
func AskManyFromOptions(message string, options, d []string) ([]string, error) {
	var choices []string
	return choices, survey.AskOne(
		&survey.MultiSelect{
			Message: message,
			Options: options,
			Default: d,
			Filter:  selectFilter,
		},
		&choices,
		nil,
	)
}

// Confirm interactively ask the user a yes/no question
func Confirm(message string, d bool) (bool, error) {
	v := false
	return v, survey.AskOne(
		&survey.Confirm{
			Message: message,
			Default: d,
		},
		&v,
		nil,
	)
}
