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

func (c *cli) Set(args ...string) error {
	if len(args) < 2 || len(args) > 3 {
		return errors.New("arguments is invalid")
	}

	switch len(args) {
	case 2:

		switch args[0] {
		case "browser":
			err := c.SetBrowser(args[1])
			if err != nil {
				return err
			}
		}

		return c.SyncConf()
	case 3:

		switch args[0] {
		case "link":
			c.Config.Links[args[1]] = args[2]

		case "search":
			c.Config.Search[args[1]] = args[2]
		}
		return c.SyncConf()

	default:
		return ErrInvalid
	}
}

func (c *cli) Del(args ...string) error {
	if len(args) != 2 {
		return errors.New("arguments is invalid")
	}

	switch args[0] {
	case "link":
		delete(c.Config.Links, args[1])
	case "search":
		delete(c.Config.Search, args[1])
	default:
		return ErrInvalid
	}
	return c.SyncConf()
}

// Open handle open browser
func (c *cli) Open(args []string, flags []string) error {
	if len(args) == 0 {
		return c.OpenDefault(flags...)
	}

	if len(args) > 1 {
		link, ok := c.Config.Search[args[0]]
		if ok {
			url, _ := utils.BuildQuery(link, "$", strings.Join(args[1:], " "))

			return browser.OpenBrowser(c.Config.Browser, append(flags, url)...)
		}

		return c.OpenSearchDefault(args[1:], flags)
	}

	link, ok := c.Config.Links[args[0]]
	if ok {
		return browser.OpenBrowser(c.Config.Browser, append(flags, link)...)
	}

	_, err := url.ParseRequestURI(args[0])
	if err != nil {
		return c.OpenSearchDefault(args, flags)
	}

	return browser.OpenBrowser(c.Config.Browser, append(flags, args[0])...)
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
func (c *cli) OpenSearchDefault(args []string, flags []string) error {
	link, ok := c.Config.Search["default"]
	if ok {
		url, _ := utils.BuildQuery(link, "$", strings.Join(args, " "))

		return browser.OpenBrowser(c.Config.Browser, append(flags, url)...)
	}

	return browser.OpenBrowser(c.Config.Browser, append(flags, strings.Join(args, " "))...)
}
