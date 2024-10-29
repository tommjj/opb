package config

import (
	"errors"
	"os/exec"

	"github.com/tommjj/go-opb/internal/interfaces"
)

var (
	ErrInvalid = errors.New("arguments is invalid")
)

type Config struct {
	// Browser app path example:"/c/br.exe"
	Browser string `json:"browser"`

	// SearchDefault name of SearchDefault
	// example:"google"
	// Search list search links
	// example:"google: https://www.google.com/search?q=$"
	Search map[string]string `json:"search,omitempty"`

	// Links is a map data
	// map[link name]url
	// example:"yt: https://www.youtube.com"
	Links map[string]string `json:"links,omitempty"`
}

type ConfigStore struct {
	Config      Config
	filestorage interfaces.IStorage
}

func New(filestorage interfaces.IStorage) *ConfigStore {
	return &ConfigStore{
		filestorage: filestorage,
		Config: Config{
			Search: map[string]string{},
			Links:  map[string]string{},
		},
	}
}

func (c *ConfigStore) Load() error {
	err := c.filestorage.Load(&c.Config)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigStore) Sync() error {
	err := c.filestorage.Sync(&c.Config)
	if err != nil {
		return err
	}

	return nil
}

// set a link and call sync
func (c *ConfigStore) SetBrowser(browserPath string) error {
	path, err := exec.LookPath(browserPath)
	if err != nil {
		return err
	}

	c.Config.Browser = path
	return c.Sync()
}

func (c *ConfigStore) Set(agr ...string) error {
	if len(agr) < 2 || len(agr) > 3 {
		return errors.New("arguments is invalid")
	}

	switch len(agr) {
	case 2:

		switch agr[0] {
		case "browser":
			err := c.SetBrowser(agr[1])
			if err != nil {
				return err
			}
		}

		return c.Sync()
	case 3:

		switch agr[0] {
		case "link":
			c.Config.Links[agr[1]] = agr[2]

		case "search":
			c.Config.Search[agr[1]] = agr[2]
		}
		return c.Sync()

	default:
		return ErrInvalid
	}
}

func (c *ConfigStore) Del(agr ...string) error {
	if len(agr) != 2 {
		return errors.New("arguments is invalid")
	}

	switch agr[0] {
	case "link":
		delete(c.Config.Links, agr[1])
	case "search":
		delete(c.Config.Search, agr[1])
	default:
		return ErrInvalid
	}
	return c.Sync()
}
