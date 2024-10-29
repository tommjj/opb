package browser

import (
	"net/url"
	"os/exec"
	"strings"

	"github.com/tommjj/go-opb/internal/config"
)

// OpenBrowser open browser
func OpenBrowser(browserPath string, arg ...string) error {
	path, err := exec.LookPath(browserPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(path, arg...)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// OpenDefault open default browser
func OpenDefault(conf config.Config, flags ...string) error {
	defaultLink, ok := conf.Links["default"]
	if ok {
		return OpenBrowser(conf.Browser, append(flags, defaultLink)...)
	}

	return OpenBrowser(conf.Browser, flags...)
}

// OpenSearchDefault open browser use default search method
func OpenSearchDefault(conf config.Config, agr []string, flags []string) error {
	link, ok := conf.Search["default"]
	if ok {
		url, _ := BuildQuery(link, "$", strings.Join(agr, " "))

		return OpenBrowser(conf.Browser, append(flags, url)...)
	}

	return OpenBrowser(conf.Browser, append(flags, strings.Join(agr, " "))...)
}

// Open handle open browser
func Open(conf config.Config, agr []string, flags []string) error {
	if len(agr) == 0 {
		return OpenDefault(conf, flags...)
	}

	if len(agr) > 1 {
		link, ok := conf.Search[agr[0]]
		if ok {
			url, _ := BuildQuery(link, "$", strings.Join(agr[1:], " "))

			return OpenBrowser(conf.Browser, append(flags, url)...)
		}

		return OpenSearchDefault(conf, agr[1:], flags)
	}

	link, ok := conf.Links[agr[0]]
	if ok {
		return OpenBrowser(conf.Browser, append(flags, link)...)
	}

	_, err := url.ParseRequestURI(agr[0])
	if err != nil {
		return OpenSearchDefault(conf, agr, flags)
	}

	return OpenBrowser(conf.Browser, append(flags, agr[0])...)
}
