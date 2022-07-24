package cmdutil

import (
	"log"
	"os"

	"github.com/lucassabreu/clockify-cli/api"
	"github.com/lucassabreu/clockify-cli/api/dto"
)

// Factory is a container/factory builder for the commands and its helpers
type Factory interface {
	Version() Version
	Config() Config
	Client() (*api.Client, error)
	GetUserID() (string, error)
	GetWorkspaceID() (string, error)
	GetWorkspace() (dto.Workspace, error)
}

type factory struct {
	version func() Version

	config func() Config
	client func() (*api.Client, error)

	getUserID      func() (string, error)
	getWorkspaceID func() (string, error)
	getWorkspace   func() (dto.Workspace, error)
}

func (f *factory) Version() Version {
	return f.version()
}

func (f *factory) Config() Config {
	return f.config()
}

func (f *factory) Client() (*api.Client, error) {
	return f.client()
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

func NewFactory(v Version) Factory {
	f := &factory{
		version: func() Version { return v },
		config:  configFunc(),
	}

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

func clientFunc(f Factory) func() (*api.Client, error) {
	var c *api.Client
	var err error

	return func() (*api.Client, error) {
		if c != nil || err != nil {
			return c, err
		}

		c, err = api.NewClient(f.Config().GetString(CONF_TOKEN))
		if f.Config().IsDebuging() {
			c.SetDebugLogger(
				log.New(os.Stdout, "DEBUG ", log.LstdFlags),
			)
		}

		return c, err
	}
}
