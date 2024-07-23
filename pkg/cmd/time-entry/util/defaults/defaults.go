package defaults

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ScanError wraps errors from scanning for the defaults file
type ScanError struct {
	Err error
}

// Error shows error message
func (s *ScanError) Error() string {
	return s.Unwrap().Error()
}

// Unwrap gives access to the error chain
func (s *ScanError) Unwrap() error {
	return s.Err
}

// DefaultsFileNotFoundErr is returned when the scan can't find any files
var DefaultsFileNotFoundErr = errors.New("defaults file not found")

const DEFAULT_FILENAME = ".clockify-defaults.json"

// DefaultTimeEntry has the default properties for the working directory
type DefaultTimeEntry struct {
	ProjectID string   `json:"project,omitempty"   yaml:"project,omitempty"`
	TaskID    string   `json:"task,omitempty"      yaml:"task,omitempty"`
	Billable  *bool    `json:"billable,omitempty"  yaml:"billable,omitempty"`
	TagIDs    []string `json:"tags,omitempty"      yaml:"tags,omitempty,flow"`
}

// ScanParam sets how ScanForDefaults should look for defaults
type ScanParam struct {
	Dir      string
	Filename string
}

// TimeEntryDefaults is a manager for the default time entry parameters on a
// folder
type TimeEntryDefaults interface {
	// Read scan the directory informed and its parents for the defaults
	// file
	Read() (DefaultTimeEntry, error)
	// Write persists the default values to the folder
	Write(DefaultTimeEntry) error
}

// NewTimeEntryDefaults creates a new instance of TimeEntryDefaults
func NewTimeEntryDefaults(p ScanParam) TimeEntryDefaults {
	return &timeEntryDefaults{
		ScanParam: p,
	}
}

type timeEntryDefaults struct {
	ScanParam
	DefaultTimeEntry
}

// FailedToOpenErr error returned when failing to open file without an explicit
// error
var FailedToOpenErr = errors.New("failed to open file")

// Write persists the default values to the folder
func (t *timeEntryDefaults) Write(d DefaultTimeEntry) error {
	println(filepath.Join(t.Dir, t.Filename))
	f, err := os.Create(filepath.Join(t.Dir, t.Filename))
	if err != nil {
		return err
	}

	if f == nil {
		return FailedToOpenErr

	}

	defer f.Close()

	if strings.HasSuffix(f.Name(), "json") {
		return json.NewEncoder(f).Encode(d)
	}

	return yaml.NewEncoder(f).Encode(d)
}

// Read scan the directory informed and its parents for the defaults
// file
func (t *timeEntryDefaults) Read() (DefaultTimeEntry, error) {
	if t.ScanParam.Filename == "" {
		t.ScanParam.Filename = DEFAULT_FILENAME
	}

	p := t.ScanParam
	dir := filepath.FromSlash(p.Dir)
	d := DefaultTimeEntry{}
	for {
		f, err := getFile(filepath.Join(dir, p.Filename))
		if err != nil {
			return d, &ScanError{
				Err: errors.Wrap(
					err, "failed to open defaults file"),
			}
		}

		if f == nil {
			nDir := filepath.Dir(dir)
			if nDir == dir {
				return d, DefaultsFileNotFoundErr
			}

			dir = nDir
			continue
		}

		if f == nil {
			return d, FailedToOpenErr

		}
		defer f.Close()

		if strings.HasSuffix(f.Name(), "json") {
			err = json.NewDecoder(f).Decode(&d)
		} else {
			err = yaml.NewDecoder(f).Decode(&d)
		}

		if err != nil {
			return d, errors.WithStack(&ScanError{
				Err: errors.Wrap(
					err, "failed to decode defaults file"),
			})
		}

		return d, nil
	}
}

func getFile(filename string) (*os.File, error) {
	stat, err := os.Stat(filepath.Join(filename))
	if err != nil || stat.IsDir() {
		return nil, nil
	}

	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, nil
	}

	return f, err
}
