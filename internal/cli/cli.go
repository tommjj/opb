package cli

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/tommjj/go-opb/internal/browser"
	"github.com/tommjj/go-opb/internal/config"
	"github.com/tommjj/go-opb/internal/interfaces"
	"github.com/tommjj/go-opb/internal/utils"
)

var (
	ErrInvalid = errors.New("arguments is invalid")
)

type cli struct {
	Config      config.Config
	filestorage interfaces.IStorage
}

func New(filestorage interfaces.IStorage) *cli {
	return &cli{
		filestorage: filestorage,
		Config: config.Config{
			Search: map[string]string{},
			Links:  map[string]string{},
		},
	}
}

func (c *cli) Reset() error {
	fmt.Println("Set your browser path")
	path := utils.GetInput()
	return c.SetBrowser(path)
}

func (c *cli) LoadConf() error {
	err := c.filestorage.Load(&c.Config)
	if err != nil {
		return err
	}

	return nil
}

func (c *cli) SyncConf() error {
	err := c.filestorage.Sync(&c.Config)
	if err != nil {
		return err
	}

	return nil
}

// set a link and call sync
func (c *cli) SetBrowser(browserPath string) error {
	path, err := exec.LookPath(browserPath)
	if err != nil {
		return err
	}

	c.Config.Browser = path
	return c.SyncConf()
}

func (c *cli) Set(agr ...string) error {
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

		return c.SyncConf()
	case 3:

		switch agr[0] {
		case "link":
			c.Config.Links[agr[1]] = agr[2]

		case "search":
			c.Config.Search[agr[1]] = agr[2]
		}
		return c.SyncConf()

	default:
		return ErrInvalid
	}
}

func (c *cli) Del(agr ...string) error {
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
	return c.SyncConf()
}

// Open handle open browser
func (c *cli) Open(agr []string, flags []string) error {
	if len(agr) == 0 {
		return c.OpenDefault(flags...)
	}

	if len(agr) > 1 {
		link, ok := c.Config.Search[agr[0]]
		if ok {
			url, _ := browser.BuildQuery(link, "$", strings.Join(agr[1:], " "))

			return browser.OpenBrowser(c.Config.Browser, append(flags, url)...)
		}

		return c.OpenSearchDefault(agr[1:], flags)
	}

	link, ok := c.Config.Links[agr[0]]
	if ok {
		return browser.OpenBrowser(c.Config.Browser, append(flags, link)...)
	}

	_, err := url.ParseRequestURI(agr[0])
	if err != nil {
		return c.OpenSearchDefault(agr, flags)
	}

	return browser.OpenBrowser(c.Config.Browser, append(flags, agr[0])...)
}

// OpenDefault open default browser
func (c *cli) OpenDefault(flags ...string) error {
	defaultLink, ok := c.Config.Links["default"]
	if ok {
		return browser.OpenBrowser(c.Config.Browser, append(flags, defaultLink)...)
	}

	return browser.OpenBrowser(c.Config.Browser, flags...)
}

// OpenSearchDefault open browser use default search method
func (c *cli) OpenSearchDefault(agr []string, flags []string) error {
	link, ok := c.Config.Search["default"]
	if ok {
		url, _ := browser.BuildQuery(link, "$", strings.Join(agr, " "))

		return browser.OpenBrowser(c.Config.Browser, append(flags, url)...)
	}

	return browser.OpenBrowser(c.Config.Browser, append(flags, strings.Join(agr, " "))...)
}
