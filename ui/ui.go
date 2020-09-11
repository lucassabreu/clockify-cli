package ui

import (
	"regexp"
	"strings"

	"github.com/lucassabreu/clockify-cli/strhlp"
	"gopkg.in/AlecAivazis/survey.v1"
)

func selectFilter(filter string, options []string) (answer []string) {
	r := strings.Join([]string{"]", "^", `\\`, "[", ".", "(", ")", "-"}, "")
	filter = regexp.MustCompile("["+r+"]+").ReplaceAllString(strhlp.Normalize(filter), "")
	filter = regexp.MustCompile(`\s+`).ReplaceAllString(filter, " ")
	filter = strings.ReplaceAll(filter, " ", ".*")
	filter = strings.ReplaceAll(filter, "*", ".*")

	regexp := regexp.MustCompile(filter)
	for _, o := range options {
		if regexp.Match([]byte(strhlp.Normalize(o))) {
			answer = append(answer, o)
		}
	}
	return answer
}

func askString(p survey.Prompt) (string, error) {
	answer := ""
	return answer, survey.AskOne(p, &answer, nil)
}

// AskForText interactively ask for one string from the user
func AskForText(message, d string) (string, error) {
	return askString(&survey.Input{
		Message: message,
		Default: d,
	})
}

// AskFromOptions interactively ask the user to choose one option or none
func AskFromOptions(message string, options []string, d string) (string, error) {
	return askString(&survey.Select{
		Message:  message,
		Options:  options,
		Default:  d,
		FilterFn: selectFilter,
	})
}

// AskManyFromOptions interactively ask the user to choose none or many option
func AskManyFromOptions(message string, options, d []string) ([]string, error) {
	var choices []string
	return choices, survey.AskOne(
		&survey.MultiSelect{
			Message:  message,
			Options:  options,
			Default:  d,
			FilterFn: selectFilter,
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
