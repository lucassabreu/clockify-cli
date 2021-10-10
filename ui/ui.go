package ui

import (
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lucassabreu/clockify-cli/strhlp"
)

func selectFilter(filter string, value string, _ int) bool {
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

// AskForText interactively ask for one string from the user
func AskForText(message, d string) (string, error) {
	return askString(&survey.Input{
		Message: message,
		Default: d,
	})
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
