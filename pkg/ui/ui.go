package ui

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/lucassabreu/clockify-cli/strhlp"
	"github.com/pkg/errors"
)

// FileReader represents the input of a terminal
type FileReader interface {
	io.Reader
	Fd() uintptr
}

// FileWriter represents the output of a terminal
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// NewUI creates a new UI instance
func NewUI(in FileReader, out FileWriter, err io.Writer) UI {
	return &ui{
		options: []survey.AskOpt{
			survey.WithStdio(in, out, err),
		},
	}
}

// UI provides functions to prompt information from a terminal
type UI interface {
	// SetPageSize changes how many entries are shown on AskFromOptions and
	// AskManyFromOptions at a time
	SetPageSize(p uint) UI
	// AskForText interactively ask for one string from the user
	AskForText(m string, opts ...InputOption) (string, error)
	// AskForValidText for a string interactively from the user and validates
	// it
	AskForValidText(
		m string, validate func(string) error, opts ...InputOption,
	) (string, error)
	// AskForDateTime interactively ask for one date and time from the user
	AskForDateTime(m, d string, ct convertTime) (time.Time, error)
	// AskForDateTimeOrNil interactively ask for one date and time from the
	// user, but allows a empty response
	AskForDateTimeOrNil(m, d string, ct convertTime) (*time.Time, error)
	// AskForInt interactively ask for one int from the user
	AskForInt(m string, d int) (int, error)
	// AskFromOptions interactively ask the user to choose one option or none
	AskFromOptions(m string, o []string, d string) (string, error)
	// AskManyFromOptions interactively ask the user to choose none or many
	// option
	AskManyFromOptions(
		m string, o, d []string, validade func([]string) error,
	) ([]string, error)
	// Confirm interactively ask the user a yes/no question
	Confirm(m string, d bool) (bool, error)
}

type ui struct {
	options []survey.AskOpt
}

func (u *ui) SetPageSize(p uint) UI {
	if p == 0 {
		p = 7
	}
	u.options = append(u.options, survey.WithPageSize(int(p)))
	return u
}

func selectFilter(filter, value string, _ int) bool {
	r := strings.Join([]string{"]", "^", `\\`, "[", ".", "(", ")", "-"}, "")
	filter = regexp.MustCompile("["+r+"]+").
		ReplaceAllString(strhlp.Normalize(filter), "")
	filter = regexp.MustCompile(`\s+`).ReplaceAllString(filter, " ")
	filter = strings.ReplaceAll(filter, " ", ".*")
	filter = strings.ReplaceAll(filter, "*", ".*")

	return regexp.MustCompile(filter).MatchString(strhlp.Normalize(value))
}

func askString(p survey.Prompt, options ...survey.AskOpt) (string, error) {
	answer := ""
	return answer, errors.WithStack(survey.AskOne(p, &answer, options...))
}

// WithSuggestion applies the suggestion function to the input question
func WithSuggestion(fn func(toComplete string) []string) InputOption {
	return func(i *survey.Input) {
		i.Suggest = fn
	}
}

// WithHelp add help to input question
func WithHelp(help string) InputOption {
	return func(i *survey.Input) {
		i.Help = help
	}
}

// WithDefault will set a default answer to the question
func WithDefault(d string) InputOption {
	return func(i *survey.Input) {
		i.Default = d
	}
}

// InputOption represets a funcion the customizes a survey.Input object
type InputOption func(*survey.Input)

// AskForValidText for a string interactively from the user and validates it
func (u *ui) AskForValidText(
	message string,
	validateFn func(string) error,
	opts ...InputOption,
) (string, error) {
	i := &survey.Input{
		Message: message,
	}

	for _, o := range opts {
		o(i)
	}

	os := u.options
	if validateFn != nil {
		os = append(os, survey.WithValidator(func(ans interface{}) error {
			return validateFn(ans.(string))
		}))
	}

	return askString(i, os...)
}

// AskForText interactively ask for one string from the user
func (u *ui) AskForText(message string, opts ...InputOption) (string, error) {
	i := &survey.Input{
		Message: message,
	}

	for _, o := range opts {
		o(i)
	}

	return askString(i, u.options...)
}

type timeAnswer struct {
	*time.Time
	convert func(string) (time.Time, error)
}

func (ans timeAnswer) validate(v interface{}) error {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}

	_, err := ans.convert(s)
	return err
}

func (ans *timeAnswer) WriteAnswer(_ string, v interface{}) error {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}

	t, err := ans.convert(s)
	if err != nil {
		return err
	}

	ans.Time = &t
	return nil
}

type convertTime func(string) (time.Time, error)

// AskForDateTime interactively ask for one date and time from the user
func (u *ui) AskForDateTime(
	name,
	value string,
	convert convertTime,
) (time.Time, error) {
	i := &survey.Input{
		Message: name + ":",
		Default: value,
	}

	t := timeAnswer{convert: convert}
	opts := make([]survey.AskOpt, 0)
	opts = append(opts, u.options...)
	opts = append(opts,
		survey.WithValidator(survey.Required),
		survey.WithValidator(t.validate),
	)

	for {
		err := survey.AskOne(i, &t, opts...)
		if err == terminal.InterruptErr {
			return time.Time{}, err
		}

		if t.Time != nil {
			return *t.Time, err
		}
	}
}

func (u *ui) AskForDateTimeOrNil(
	name,
	value string,
	convert convertTime,
) (*time.Time, error) {
	t := timeAnswer{convert: convert}
	opts := []survey.AskOpt{survey.WithValidator(t.validate)}
	opts = append(opts, u.options...)
	return t.Time, survey.AskOne(
		&survey.Input{
			Message: name + " (leave it blank for empty):",
			Default: value,
		},
		&t,
		opts...,
	)
}

// AskForInt interactively ask for one int from the user
func (u *ui) AskForInt(message string, d int) (int, error) {
	opts := []survey.AskOpt{survey.WithValidator(func(ans interface{}) error {
		v, ok := ans.(string)
		if !ok {
			return fmt.Errorf("needs to be a string")
		}

		_, err := strconv.Atoi(v)
		return err
	})}
	opts = append(opts, u.options...)
	return d, survey.AskOne(
		&survey.Input{
			Message: message,
			Default: strconv.Itoa(d),
		},
		&d,
		opts...,
	)
}

// AskFromOptions interactively ask the user to choose one option or none
func (u *ui) AskFromOptions(message string, options []string, d string) (string, error) {
	p := &survey.Select{
		Message: message,
		Options: options,
		Filter:  selectFilter,
	}

	if d != "" && strhlp.Search(d, options) != -1 {
		p.Default = d
	}

	return askString(p, u.options...)
}

// AskManyFromOptions interactively ask the user to choose none or many option
func (u *ui) AskManyFromOptions(
	message string,
	opts, d []string,
	validateFn func([]string) error,
) ([]string, error) {
	var choices []string

	os := u.options
	if validateFn != nil {
		os = append(os, survey.WithValidator(func(ans interface{}) error {
			o := ans.([]survey.OptionAnswer)
			s := make([]string, len(o))
			for i := range o {
				s[i] = o[i].Value
			}
			return validateFn(s)
		}))
	}

	return choices, survey.AskOne(
		&survey.MultiSelect{
			Message: message,
			Options: opts,
			Default: d,
			Filter:  selectFilter,
		},
		&choices,
		os...,
	)
}

// Confirm interactively ask the user a yes/no question
func (u *ui) Confirm(message string, d bool) (bool, error) {
	v := false
	return v, survey.AskOne(
		&survey.Confirm{
			Message: message,
			Default: d,
		},
		&v,
		u.options...,
	)
}
