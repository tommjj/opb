package browser

import (
	"os/exec"
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
