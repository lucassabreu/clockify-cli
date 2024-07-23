package cmdutil

import (
	"log"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
	"github.com/lucassabreu/clockify-cli/pkg/cmd/time-entry/util/defaults"
	"github.com/lucassabreu/clockify-cli/pkg/ui"
)

// Factory is a container/factory builder for the commands and its helpers
type Factory interface {
	// Version of the CLI
	Version() Version

	// Config returns configurations set by the user
	Config() Config
	// Client builds a client for Clockify's API
	Client() (api.Client, error)
	// UI builds a control to prompt information from the user
	UI() ui.UI
	// TimeEntryDefaults manages the default properties of a time entry
	TimeEntryDefaults() defaults.TimeEntryDefaults

	// GetUserID returns the current user id
	GetUserID() (string, error)
	// GetWorkspaceID returns the current workspace id
	GetWorkspaceID() (string, error)
	// GetWorkspaceID returns the current workspace
	GetWorkspace() (dto.Workspace, error)
}

type factory struct {
	version func() Version

	config            func() Config
	client            func() (api.Client, error)
	ui                func() ui.UI
	timeEntryDefaults func() defaults.TimeEntryDefaults

	getUserID      func() (string, error)
	getWorkspaceID func() (string, error)
	getWorkspace   func() (dto.Workspace, error)
}

// TimeEntryDefaults manages the default properties of a time entry
func (f *factory) TimeEntryDefaults() defaults.TimeEntryDefaults {
	return f.timeEntryDefaults()
}

func (f *factory) Version() Version {
	return f.version()
}

func (f *factory) Config() Config {
	return f.config()
}

func (f *factory) Client() (api.Client, error) {
	return f.client()
}

func (f *factory) UI() ui.UI {
	return f.ui()
}

func (f *factory) GetUserID() (string, error) {
	return f.getUserID()
}

func (f *factory) GetWorkspaceID() (string, error) {
	return f.getWorkspaceID()
}

func (f *factory) GetWorkspace() (dto.Workspace, error) {
	return f.getWorkspace()
}

// NewFactory creates a new instance of Factory
func NewFactory(v Version) Factory {
	f := &factory{
		version:           func() Version { return v },
		config:            configFunc(),
		timeEntryDefaults: getTED(),
	}

	f.ui = getUi(f)

	f.client = clientFunc(f)

	f.getUserID = getUserIDFunc(f)

	f.getWorkspace = getWorkspaceFunc(f)
	f.getWorkspaceID = getWorkspaceIDFunc(f)

	return f
}

func getUserIDFunc(f Factory) func() (string, error) {
	var userID string
	var err error
	return func() (string, error) {
		if userID != "" || err != nil {
			return userID, err
		}

		userID = f.Config().GetString(CONF_USER_ID)
		if userID != "" {
			return userID, err
		}

		client, err := f.Client()
		if err != nil {
			return userID, err
		}

		u, err := client.GetMe()
		if err != nil {
			return userID, err
		}

		userID = u.ID
		return userID, err
	}
}

func getWorkspaceFunc(f Factory) func() (dto.Workspace, error) {
	var w *dto.Workspace
	var err error
	return func() (dto.Workspace, error) {
		if w != nil {
			return *w, nil
		}

		if err != nil {
			return dto.Workspace{}, err
		}

		id, err := f.GetWorkspaceID()
		if err != nil {
			return dto.Workspace{}, err
		}

		client, err := f.Client()
		if err != nil {
			return dto.Workspace{}, err
		}

		oW, err := client.GetWorkspace(api.GetWorkspace{ID: id})
		if err != nil {
			return dto.Workspace{}, err
		}

		w = &oW
		return *w, err
	}
}

func getWorkspaceIDFunc(f Factory) func() (string, error) {
	var w string
	var err error
	return func() (string, error) {
		if w != "" || err != nil {
			return w, err
		}

		w = f.Config().GetString(CONF_WORKSPACE)
		if w != "" {
			return w, err
		}

		client, err := f.Client()
		if err != nil {
			return w, err
		}

		u, err := client.GetMe()
		w = u.DefaultWorkspace

		return w, err
	}
}

func clientFunc(f Factory) func() (api.Client, error) {
	var c api.Client
	var err error

	return func() (api.Client, error) {
		if c != nil || err != nil {
			return c, err
		}

		c, err = api.NewClient(f.Config().GetString(CONF_TOKEN))
		if err != nil {
			return c, err
		}

		ll := f.Config().LogLevel()
		if ll == LOG_LEVEL_NONE {
			return c, err
		}

		c.SetInfoLogger(
			log.New(os.Stdout, "INFO  ", log.LstdFlags),
		)

		if ll == LOG_LEVEL_INFO {
			return c, err
		}

		c.SetDebugLogger(
			log.New(os.Stdout, "DEBUG ", log.LstdFlags),
		)

		return c, err
	}
}

func getUi(f Factory) func() ui.UI {
	var i ui.UI
	return func() ui.UI {
		if i == nil {
			i = ui.NewUI(os.Stdin, os.Stdout, os.Stderr)
			i.SetPageSize(uint(f.Config().InteractivePageSize()))
		}

		return i
	}

}

func getTED() func() defaults.TimeEntryDefaults {
	var ted defaults.TimeEntryDefaults
	return func() defaults.TimeEntryDefaults {
		if ted != nil {
			return ted
		}

		wd, _ := os.Getwd()
		ted = defaults.NewTimeEntryDefaults(defaults.ScanParam{
			Dir: wd,
		})

		return ted
	}
}
